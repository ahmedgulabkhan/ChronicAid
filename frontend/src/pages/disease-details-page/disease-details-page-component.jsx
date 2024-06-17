import React, { useState, useEffect } from "react";
import { useParams, Link } from 'react-router-dom';
import { useSelector } from "react-redux";

import './disease-details-page-styles.scss';

const DiseaseDetailsPage = () => {
    const username = useSelector(state => state.user.username);
    const isLoggedIn = useSelector(state => state.user.isLoggedIn);
    const { diseaseName } = useParams();
    const [metricsForDisease, setMetricsForDisease] = useState([]);

    useEffect(() => {
        if (username === "") return;
        fetch('http://localhost:8080/api/diseases/metrics/all/' +  diseaseName, {
            method: 'GET',
            headers: {
                'Accept': 'application/json',
                'Content-Type': 'application/json'
            },
            credentials: 'include',
        })
        .then(res => res.json())
        .then(json => setMetricsForDisease(json))
        .catch(err => console.log(err));
    }, [username]);

    const toUpperCamelCase = (str) => {
        return str
            .split('-')
            .map(word => word.charAt(0).toUpperCase() + word.slice(1).toLowerCase())
            .join(' ');
    }

    return (
        (isLoggedIn === true) ?
        <div className='disease-details-page'>
            <h1>Metrics for {toUpperCamelCase(diseaseName)}</h1>

            <div className="disease-metrics-container-wrapper">
            {
                    (() => {
                        let metrics = metricsForDisease;
                        let numRows = metrics.length/3;
                        if (numRows%4 != 0) numRows += 1;
                        const elements = [];
                        let startIndex = 0;
                        let endIndex = 3;

                        for (let i=0; i<numRows; i++) {
                            let rowMetrics = metrics.splice(startIndex, endIndex);
                            elements.push(
                                <div className="disease-metrics-container-row">
                                    {rowMetrics.map((item, index) => (
                                        <Link to={`/diseases/${diseaseName}/${item}`} className="diseases-card-link">
                                            <div className="disease-metrics-card">
                                                <strong>{item}</strong>
                                            </div>
                                        </Link>
                                    ))}
                                </div>
                            );
                        }
                        return elements;
                    })()
                }
            </div>
        </div> :
        <div className='disease-details-page'>
            <h1>Metrics for {toUpperCamelCase(diseaseName)}</h1>
            <p>This is empty. Login to know more details</p>
        </div>
    );
}

export default DiseaseDetailsPage;