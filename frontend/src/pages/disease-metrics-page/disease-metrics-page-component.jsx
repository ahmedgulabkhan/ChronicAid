import React, { useState } from "react";
import { useEffect } from "react";
import { useSelector } from "react-redux";
import { useParams, Link } from 'react-router-dom';

import './disease-metrics-page-styles.scss';

const DiseaseMetricsPage = () => {
    const username = useSelector(state => state.user.username);
    const isLoggedIn = useSelector(state => state.user.isLoggedIn);
    const [userDiseaseMetrics, setUserDiseaseMetrics] = useState([]);
    const [selectedOption, setSelectedOption] = useState('last7days');
    const { diseaseName } = useParams();
    const { metricName } = useParams();
    const Swal = require('sweetalert2')

    const onShowAlert = (data) => {
        let title;
        let text;
        Object.entries(data).map(([key, value]) => {
            title = value.split('-')[0];
            text = "The value for this metric was captured on " + convertEpochDayToDate(key) + " and the status is " + value.split('-')[1];
        });
        Swal.fire({
            title: title,
            text: text
        });
    }

    const convertMapToArray = (mapData) => {
        let array = []
        Object.entries(mapData).map(([key, item]) => {
            array.push({[key]: item})
        });
        return array;
    }

    const convertEpochDayToDate = (epochDay) => {
        const milliseconds = epochDay * 86400000;
        const date = new Date(milliseconds);
        const options = { year: 'numeric', month: 'long', day: 'numeric' };
        return date.toLocaleDateString('en-US', options);
    } 

    const handleOptionChange = (option) => {
        setSelectedOption(option);
    }

    useEffect(() => {
        if (username === "") return;
        fetch('http://localhost:8080/api/diseases/metrics/' + username + '/' + diseaseName + '/' + selectedOption, {
            method: 'GET',
            headers: {
                'Accept': 'application/json',
                'Content-Type': 'application/json'
            },
            credentials: 'include',
        })
        .then(res => res.json())
        .then(json => setUserDiseaseMetrics(convertMapToArray(json)))
        .catch(err => console.log(err));
    }, [username, selectedOption]);

    return (
        (isLoggedIn === true) ?
        <div className="disease-metrics-page">
            <h1>{metricName}</h1>
            <div className="toggle-switch">
                <button
                className={selectedOption === 'last7days' ? 'active' : ''}
                onClick={() => handleOptionChange('last7days')}
                >
                    Last 7 days
                </button>
                <button
                className={selectedOption === 'last28days' ? 'active' : ''}
                onClick={() => handleOptionChange('last28days')}
                >
                    Last 28 days
                </button>
            </div>

            <div className="disease-metrics-container-wrapper">
                {
                    (() => {                        
                        if (userDiseaseMetrics.length == 0) return;
                        let metrics = userDiseaseMetrics;
                        let metricsForMetricName = [];

                        metrics.map(item => {
                            Object.entries(item).map(([key, value]) => {
                                if (key === metricName) {
                                    metricsForMetricName = convertMapToArray(value);
                                }
                            });
                        })

                        let numRows = metricsForMetricName.length/4;
                        if (numRows%4 != 0) numRows += 1;
                        const elements = [];
                        let startIndex = 0;
                        let endIndex = 4;

                        for (let i=0; i<numRows; i++) {
                            let rowMetrics = metricsForMetricName.splice(startIndex, endIndex);
                            elements.push(
                                <div className="disease-metrics-container-row">
                                    {rowMetrics.map((item, index) => (
                                        <div className="disease-metrics-card" onClick={() => onShowAlert(item)}>
                                            {Object.entries(item).map(([key, value]) => (
                                                <p key={key}>
                                                    <strong>{convertEpochDayToDate(key)}</strong><br />
                                                    <strong>Value:</strong> {value.split('-')[0]}<br />
                                                    <strong>Status:</strong> {value.split('-')[1]}
                                                </p>
                                            ))}
                                        </div>
                                    ))}
                                </div>
                            );
                        }
                        return elements;
                    })()
                }
            </div>
        </div> :
        <div className='disease-metrics-page'>
            <h1>{metricName}</h1>
            <p>This is empty. Login to know more details</p>
        </div>
    );
}

export default DiseaseMetricsPage;