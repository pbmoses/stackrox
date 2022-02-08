import React from 'react';
import {
    Badge,
    Button,
    ButtonVariant,
    Divider,
    Flex,
    FlexItem,
    PageSection,
    PageSectionVariants,
    Title,
} from '@patternfly/react-core';
import { TableComposable, Thead, Tbody, Tr, Th, Td } from '@patternfly/react-table';
import { useParams, Link } from 'react-router-dom';
import resolvePath from 'object-resolve-path';
import pluralize from 'pluralize';

import ACSEmptyState from 'Components/ACSEmptyState';
import ButtonLink from 'Components/PatternFly/ButtonLink';
import useTableSelection from 'hooks/useTableSelection';
import useIntegrationPermissions from '../hooks/useIntegrationPermissions';
import usePageState from '../hooks/usePageState';
import { Integration, getIsAPIToken, getIsClusterInitBundle } from '../utils/integrationUtils';
import tableColumnDescriptor from '../utils/tableColumnDescriptor';
import DownloadCAConfigBundle from '../ClusterInitBundles/DownloadCAConfigBundle';

type TableCellProps = {
    row: Integration;
    column: {
        Header: string;
        accessor: (data) => string | string;
    };
};

function TableCellValue({ row, column }: TableCellProps): React.ReactElement {
    let value: string;
    if (typeof column.accessor === 'function') {
        value = column.accessor(row).toString();
    } else {
        value = resolvePath(row, column.accessor).toString();
    }
    return <div>{value || '-'}</div>;
}

function getNewButtonText(type) {
    if (type === 'apitoken') {
        return 'Generate token';
    }
    if (type === 'clusterInitBundle') {
        return 'Generate bundle';
    }
    return 'New integration';
}

type IntegrationsTableProps = {
    title: string;
    integrations: Integration[];
    hasMultipleDelete: boolean;
    onDeleteIntegrations: (integration) => void;
};

function IntegrationsTable({
    title,
    integrations,
    hasMultipleDelete,
    onDeleteIntegrations,
}: IntegrationsTableProps): React.ReactElement {
    const permissions = useIntegrationPermissions();
    const { source, type } = useParams();
    const { getPathToCreate, getPathToEdit, getPathToViewDetails } = usePageState();
    const columns = [...tableColumnDescriptor[source][type]];
    const {
        selected,
        allRowsSelected,
        numSelected,
        hasSelections,
        onSelect,
        onSelectAll,
        getSelectedIds,
    } = useTableSelection<Integration>(integrations);

    const isAPIToken = getIsAPIToken(source, type);
    const isClusterInitBundle = getIsClusterInitBundle(source, type);

    function onDeleteIntegrationHandler() {
        const ids = getSelectedIds();
        onDeleteIntegrations(ids);
    }

    const newButtonText = getNewButtonText(type);

    return (
        <>
            <Flex className="pf-u-p-md">
                <FlexItem
                    spacer={{ default: 'spacerMd' }}
                    alignSelf={{ default: 'alignSelfCenter' }}
                >
                    <Flex alignSelf={{ default: 'alignSelfCenter' }}>
                        <FlexItem
                            spacer={{ default: 'spacerMd' }}
                            alignSelf={{ default: 'alignSelfCenter' }}
                        >
                            <Title headingLevel="h2" className="pf-u-color-100 pf-u-ml-sm">
                                {title}
                            </Title>
                        </FlexItem>
                        <FlexItem
                            spacer={{ default: 'spacerMd' }}
                            alignSelf={{ default: 'alignSelfCenter' }}
                        >
                            <Badge isRead>{integrations.length}</Badge>
                        </FlexItem>
                    </Flex>
                </FlexItem>
                <FlexItem align={{ default: 'alignRight' }}>
                    <Flex>
                        {hasSelections && hasMultipleDelete && permissions[source].write && (
                            <FlexItem spacer={{ default: 'spacerMd' }}>
                                <Button variant="danger" onClick={onDeleteIntegrationHandler}>
                                    Delete {numSelected} selected{' '}
                                    {pluralize('integration', numSelected)}
                                </Button>
                            </FlexItem>
                        )}
                        {isClusterInitBundle && (
                            <FlexItem spacer={{ default: 'spacerMd' }}>
                                <DownloadCAConfigBundle />
                            </FlexItem>
                        )}
                        {permissions[source].write && (
                            <FlexItem spacer={{ default: 'spacerMd' }}>
                                <ButtonLink
                                    to={getPathToCreate(source, type)}
                                    variant={ButtonVariant.primary}
                                    data-testid="add-integration"
                                >
                                    {newButtonText}
                                </ButtonLink>
                            </FlexItem>
                        )}
                    </Flex>
                </FlexItem>
            </Flex>
            <Divider component="div" />
            <PageSection
                isFilled
                padding={{ default: 'noPadding' }}
                hasOverflowScroll
                variant={PageSectionVariants.light}
            >
                {integrations.length > 0 ? (
                    <TableComposable variant="compact" isStickyHeader>
                        <Thead>
                            <Tr>
                                {hasMultipleDelete && (
                                    <Th
                                        select={{
                                            onSelect: onSelectAll,
                                            isSelected: allRowsSelected,
                                        }}
                                    />
                                )}
                                {columns.map((column) => {
                                    return (
                                        <Th key={column.Header} modifier="wrap">
                                            {column.Header}
                                        </Th>
                                    );
                                })}
                                <Th aria-label="Row actions" />
                            </Tr>
                        </Thead>
                        <Tbody>
                            {integrations.map((integration, rowIndex) => {
                                const { id } = integration;
                                const actionItems = [
                                    {
                                        title: (
                                            <Link to={getPathToEdit(source, type, id)}>
                                                Edit integration
                                            </Link>
                                        ),
                                        isHidden: isAPIToken || isClusterInitBundle,
                                    },
                                    {
                                        title: (
                                            <div className="pf-u-danger-color-100">
                                                Delete Integration
                                            </div>
                                        ),
                                        onClick: () => onDeleteIntegrations([integration.id]),
                                    },
                                ].filter((actionItem) => {
                                    return !actionItem?.isHidden;
                                });
                                return (
                                    <Tr key={integration.id}>
                                        {hasMultipleDelete && (
                                            <Td
                                                key={integration.id}
                                                select={{
                                                    rowIndex,
                                                    onSelect,
                                                    isSelected: selected[rowIndex],
                                                }}
                                            />
                                        )}
                                        {columns.map((column) => {
                                            if (column.Header === 'Name') {
                                                return (
                                                    <Td key="name">
                                                        <Button
                                                            variant={ButtonVariant.link}
                                                            isInline
                                                            component={(props) => (
                                                                <Link
                                                                    {...props}
                                                                    to={getPathToViewDetails(
                                                                        source,
                                                                        type,
                                                                        id
                                                                    )}
                                                                />
                                                            )}
                                                        >
                                                            <TableCellValue
                                                                row={integration}
                                                                column={column}
                                                            />
                                                        </Button>
                                                    </Td>
                                                );
                                            }
                                            return (
                                                <Td key={column.Header}>
                                                    <TableCellValue
                                                        row={integration}
                                                        column={column}
                                                    />
                                                </Td>
                                            );
                                        })}
                                        <Td
                                            actions={{
                                                items: actionItems,
                                                disable: !permissions[source].write,
                                            }}
                                            className="pf-u-text-align-right"
                                        />
                                    </Tr>
                                );
                            })}
                        </Tbody>
                    </TableComposable>
                ) : (
                    <ACSEmptyState
                        key="no-results"
                        title="No integrations of this type are currently configured."
                    />
                )}
            </PageSection>
        </>
    );
}

export default IntegrationsTable;
