import React from "react";

import './homepage-styles.scss';
import TrackProgress from "../../components/track-progress/track-progress-component";
import YourDiseases from "../../components/your-diseases/your-diseases-component";

const HomePage = () => {
    return (
        <div className='home-page'>
            <TrackProgress showAllMetrics={false} />
            <YourDiseases showAllDiseases={false} />
        </div>
    );
}

export default HomePage;