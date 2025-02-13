import { QueryValue } from 'hooks/useURLParameter';

type SimulationOn = {
    isOn: true;
    type: 'baseline' | 'networkPolicy';
};

type SimulationOff = {
    isOn: false;
};

export type Simulation = SimulationOn | SimulationOff;

function getSimulation(simulationQueryValue: QueryValue): Simulation {
    if (
        !simulationQueryValue ||
        (simulationQueryValue !== 'baseline' && simulationQueryValue !== 'networkPolicy')
    ) {
        const simulation: Simulation = {
            isOn: false,
        };
        return simulation;
    }
    const simulation: Simulation = {
        isOn: true,
        type: simulationQueryValue,
    };
    return simulation;
}

export default getSimulation;
