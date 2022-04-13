import React, { useEffect, useMemo, useState, ReactElement } from 'react';
import { PageSection, Bullseye, Alert, Divider, Title } from '@patternfly/react-core';

import { fetchAlerts, fetchAlertCount, CanceledPromiseError } from 'services/AlertsService';
import { getSearchOptionsForCategory } from 'services/SearchService';

import useEntitiesByIdsCache from 'hooks/useEntitiesByIdsCache';
import LIFECYCLE_STAGES from 'constants/lifecycleStages';
import VIOLATION_STATES from 'constants/violationStates';
import { ENFORCEMENT_ACTIONS } from 'constants/enforcementActions';
import { SEARCH_CATEGORIES } from 'constants/searchOptions';

import useEffectAfterFirstRender from 'hooks/useEffectAfterFirstRender';
import useURLSort, { SortOption } from 'hooks/patternfly/useURLSort';
import useURLSearch from 'hooks/useURLSearch';
import useURLPagination from 'hooks/useURLPagination';
import { checkForPermissionErrorMessage } from 'utils/permissionUtils';
import SearchFilterInput from 'Components/SearchFilterInput';
import ViolationsTablePanel from './ViolationsTablePanel';
import tableColumnDescriptor from './violationTableColumnDescriptors';

import './ViolationsTablePage.css';

const runAfter5Seconds = (fn: () => void) => setTimeout(fn, 5000);

const searchCategory = SEARCH_CATEGORIES.ALERTS;

function ViolationsTablePage(): ReactElement {
    // Handle changes to applied search options.
    const [searchOptions, setSearchOptions] = useState<string[]>([]);
    const { searchFilter, setSearchFilter } = useURLSearch();

    const hasExecutableFilter =
        Object.keys(searchFilter).length &&
        Object.values(searchFilter).some((filter) => filter !== '');

    const [isViewFiltered, setIsViewFiltered] = useState(hasExecutableFilter);

    // Handle changes in the current table page.
    const { page, perPage, setPage, setPerPage } = useURLPagination(50);

    // Handle changes in the currently displayed violations.
    const [currentPageAlerts, setCurrentPageAlerts] = useEntitiesByIdsCache();
    const [currentPageAlertsErrorMessage, setCurrentPageAlertsErrorMessage] = useState('');
    const [alertCount, setAlertCount] = useState(0);

    // To handle page/count refreshing.
    const [pollEpoch, setPollEpoch] = useState(0);

    // To handle sort options.
    const columns = tableColumnDescriptor;
    const sortFields = useMemo(
        () => columns.flatMap(({ sortField }) => (sortField ? [sortField] : [])),
        [columns]
    );

    const defaultSortOption: SortOption = {
        field: 'Violation Time',
        direction: 'desc',
    };
    const { sortOption, getSortParams } = useURLSort({
        sortFields,
        defaultSortOption,
    });

    useEffectAfterFirstRender(() => {
        if (hasExecutableFilter && !isViewFiltered) {
            // If the user applies a filter to a previously unfiltered table, return to page 1
            setIsViewFiltered(true);
            setPage(1);
        } else if (!hasExecutableFilter && isViewFiltered) {
            // If the user clears all filters after having previously applied filters, return to page 1
            setIsViewFiltered(false);
            setPage(1);
        }
    }, [hasExecutableFilter, isViewFiltered, setIsViewFiltered, setPage]);

    useEffectAfterFirstRender(() => {
        // Prevent viewing a page beyond the maximum page count
        if (page > Math.ceil(alertCount / perPage)) {
            setPage(1);
        }
    }, [alertCount, perPage, setPage]);

    // When any of the deps to this effect change, we want to reload the alerts and count.
    useEffect(() => {
        // TODO Only clear error message when a req is not in progress?
        setCurrentPageAlertsErrorMessage('');

        const alertRequest = fetchAlerts(searchFilter, sortOption, page - 1, perPage);
        const countRequest = fetchAlertCount(searchFilter);

        Promise.all([alertRequest, countRequest])
            .then(([alerts, counts]) => {
                setCurrentPageAlerts(alerts);
                setAlertCount(counts);
            })
            .catch((error) => {
                if (error instanceof CanceledPromiseError) {
                    return;
                }
                setCurrentPageAlerts([]);
                const parsedMessage = checkForPermissionErrorMessage(error);
                setCurrentPageAlertsErrorMessage(parsedMessage);
            });

        const timeoutId = runAfter5Seconds(() => {
            setPollEpoch(pollEpoch + 1);
        });

        return () => {
            clearTimeout(timeoutId);
            alertRequest.cancel();
            countRequest.cancel();
        };
    }, [
        searchFilter,
        page,
        sortOption,
        pollEpoch,
        setCurrentPageAlerts,
        setCurrentPageAlertsErrorMessage,
        setAlertCount,
        perPage,
    ]);

    useEffect(() => {
        getSearchOptionsForCategory(searchCategory)
            .then(setSearchOptions)
            .catch(() => {
                // Is there a reasonable way to recover from a possible error here?
                // Right now, ignoring this error simply disables the search filter.
            });
    }, [setSearchOptions]);

    // We need to be able to identify which alerts are runtime or attempted, and which are not by id.
    const resolvableAlerts: Set<string> = new Set(
        currentPageAlerts
            .filter(
                (alert) =>
                    alert.lifecycleStage === LIFECYCLE_STAGES.RUNTIME ||
                    alert.state === VIOLATION_STATES.ATTEMPTED
            )
            .map((alert) => alert.id as string)
    );

    const excludableAlerts = currentPageAlerts.filter(
        (alert) =>
            alert.enforcementAction !== ENFORCEMENT_ACTIONS.FAIL_DEPLOYMENT_CREATE_ENFORCEMENT
    );

    return (
        <>
            <PageSection variant="light" id="violations-table">
                <Title headingLevel="h1">Violations</Title>
                <Divider className="pf-u-py-md" />
                <SearchFilterInput
                    className="theme-light"
                    handleChangeSearchFilter={setSearchFilter}
                    placeholder="Filter violations"
                    searchCategory={searchCategory}
                    searchFilter={searchFilter}
                    searchOptions={searchOptions}
                />
            </PageSection>
            <PageSection variant="default">
                {currentPageAlertsErrorMessage ? (
                    <Bullseye>
                        <Alert variant="danger" title={currentPageAlertsErrorMessage} />
                    </Bullseye>
                ) : (
                    <PageSection variant="light">
                        <ViolationsTablePanel
                            violations={currentPageAlerts}
                            violationsCount={alertCount}
                            currentPage={page}
                            setCurrentPage={setPage}
                            resolvableAlerts={resolvableAlerts}
                            excludableAlerts={excludableAlerts}
                            perPage={perPage}
                            setPerPage={setPerPage}
                            getSortParams={getSortParams}
                            columns={columns}
                        />
                    </PageSection>
                )}
            </PageSection>
        </>
    );
}

export default ViolationsTablePage;
