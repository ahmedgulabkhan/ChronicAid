import React from "react";
import './diseases-page-styles.scss';
import YourDiseases from "../../components/your-diseases/your-diseases-component";

const DiseasesPage = () => {
    return (
        <div className="diseases-page">
            <YourDiseases showAllDiseases={true} />
        </div>
    );
}

export default DiseasesPage;