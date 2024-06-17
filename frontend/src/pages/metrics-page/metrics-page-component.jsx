import React from "react";

import TrackProgress from "../../components/track-progress/track-progress-component";

import './metrics-page-styles.scss';

const MetricsPage = () => {
    return (
        <div className="metrics-page">
            <TrackProgress showAllMetrics={true} />
        </div>
    );
}

export default MetricsPage;