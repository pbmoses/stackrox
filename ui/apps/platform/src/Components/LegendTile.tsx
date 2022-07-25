import React from 'react';
import { ListItem } from '@patternfly/react-core';
import deployment from 'images/legend-icons/deployment.svg';
import deploymentActiveConnection from 'images/legend-icons/deployment-active-connection.svg';
import deploymentAllowedConnection from 'images/legend-icons/deployment-allowed-connection.svg';
import deploymentAllowedConnections from 'images/legend-icons/deployment-allowed-connections.svg';
import nonIsolatedDeploymentAllowed from 'images/legend-icons/non-isolated-deployment-allowed.svg';
import deploymentExternalConnections from 'images/legend-icons/deployment-with-external-flows.svg';
import disallowedConnection from 'images/legend-icons/disallowed-connection.svg';
import disallowedDeployment from 'images/legend-icons/disallowed-deployment.svg';
import namespace from 'images/legend-icons/namespace.svg';
import namespaceAllowed from 'images/legend-icons/namespace-allowed.svg';
import namespaceConnection from 'images/legend-icons/namespace-connection.svg';
import namespaceDisallowedConnection from 'images/legend-icons/namespace-disallowed-connection.svg';
import namespaceEgressIngress from 'images/legend-icons/namespace-egress-ingress.svg';

const svgMapping = {
    'active-connection': deploymentActiveConnection,
    'allowed-connection': deploymentAllowedConnection,
    namespace,
    deployment,
    'disallowed-connection': disallowedConnection,
    'disallowed-deployment': disallowedDeployment,
    'namespace-allowed-connection': namespaceAllowed,
    'namespace-connection': namespaceConnection,
    'non-isolated-deployment-allowed': nonIsolatedDeploymentAllowed,
    'namespace-disallowed-connection': namespaceDisallowedConnection,
    'namespace-egress-ingress': namespaceEgressIngress,
    'deployment-external-connections': deploymentExternalConnections,
    'deployment-allowed-connections': deploymentAllowedConnections,
};

type LegendTileProps = {
    name: string;
    description: string;
};

const LegendTile = ({ name, description }: LegendTileProps) => (
    <ListItem className="pf-u-font-size-xs pf-u-display-flex pf-u-align-items-center">
        <img src={svgMapping[name]} alt={name} style={{ width: '18px', height: '18px' }} />
        <span className="pf-u-ml-sm">{description}</span>
    </ListItem>
);

export default LegendTile;
