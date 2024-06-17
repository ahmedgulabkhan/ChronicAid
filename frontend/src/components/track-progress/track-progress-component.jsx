import React, { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { useSelector } from "react-redux";

import './track-progress-styles.scss';

const TrackProgress = (props) => {
    const username = useSelector(state => state.user.username)
    const isLoggedIn = useSelector(state => state.user.isLoggedIn);
    const [userMetricsForToday, setUserMetricsForToday] = useState([]);
    const Swal = require('sweetalert2');
    
    const onShowAlert = (data, metricName) => {
        let title;
        let text;
        Object.entries(data).map(([key, value]) => {
            title = value.split('-')[0];
            text = "The above is the metric value captured for " + metricName + " on " + convertEpochDayToDate(key) + " and the status is " + value.split('-')[1];
        });
        Swal.fire({
            title: title,
            text: text
        });
    }

    const convertMapToArray = (mapData) => {
        let array = []
        Object.entries(mapData).map(([key, item]) => array.push({[key]: item}));
        return array;
    }

    const convertEpochDayToDate = (epochDay) => {
        const milliseconds = epochDay * 86400000;
        const date = new Date(milliseconds);
        const options = { year: 'numeric', month: 'long', day: 'numeric' };
        return date.toLocaleDateString('en-US', options);
    }

    useEffect(() => {
        if (username === undefined || username === "") return;
        fetch('http://localhost:8080/api/diseases/metrics/alllast7days/' + username, {
            method: 'GET',
            headers: {
                'Accept': 'application/json',
                'Content-Type': 'application/json'
            },
            credentials: 'include',
        })
        .then(res => res.json())
        .then(json => {
            console.log(json)
            setUserMetricsForToday(convertMapToArray(json));
        })
        .catch(err => console.log(err));
    }, [username]);

    const handleAddNewMetric = () => {
        Swal.fire({
            title: 'Add New Metric',
            html: `
                <input type="text" id="swal-input1" class="swal2-input" placeholder="Disease Name">
                <input type="text" id="swal-input2" class="swal2-input" placeholder="Metric Name">
                <input type="text" id="swal-input3" class="swal2-input" placeholder="Metric Value">
                `,
            focusConfirm: false,
            showCancelButton: true,
            confirmButtonText: 'Submit',
            cancelButtonText: 'Cancel',
            preConfirm: () => {
                return {
                    disease_name: document.getElementById('swal-input1').value,
                    metric_name: document.getElementById('swal-input2').value,
                    metric_value: document.getElementById('swal-input3').value
                };
            }
        }).then((result) => {
            if (result.isConfirmed && result.value !== "") {
                const currentEpochDay = Math.floor(Date.now() / 86400000)
                const data = {
                    metric_name: result.value.metric_name,
                    metric_values: {
                        [currentEpochDay]: result.value.metric_value
                    }
                }
                fetch('http://localhost:8080/api/diseases/metrics/' + username + "/" + result.value.disease_name, {
                    method: 'POST',
                    headers: {
                        'Accept': 'application/json',
                        'Content-Type': 'application/json'
                    },
                    credentials: 'include',
                    body: JSON.stringify(data)
                })
                .then(res => res.json())
                .then(json => {
                    Swal.fire({
                        title: "AI Advice",
                        text: json.message
                    });
                })
                .catch(err => console.log(err));
            }
        });
    }

    return (
        (isLoggedIn === true) ?
        <div className='track-progress'>
            {
                (!props.showAllMetrics) ?
                <div className="heading-section">
                    <h1>Track Your Progress</h1>
                    <button className="add-new-button" onClick={handleAddNewMetric}>Add New Metric</button>
                </div> :
                <div className="heading-section">
                    <h1>Metrics</h1>
                    <button className="add-new-button">Add New Metric</button>
                </div>
            }
            <div className="disease-metrics-container-wrapper">
                {
                    (() => {
                        if (userMetricsForToday.length == 0) return;
                        let cardItemCount = 0;
                        let cardItems = [];
                        for (let i=0; i<userMetricsForToday.length; i++) {
                            Object.entries(userMetricsForToday[i]).map(([key, value]) => {
                                Object.entries(value).map(([keyInner, valueInner]) => {
                                    cardItems.push({[key]: {[keyInner]: valueInner}});
                                    cardItemCount++;
                                })
                            })
                        }

                        let cardItemsToInclude = cardItemCount;

                        if (!props.showAllMetrics) {
                            cardItemsToInclude = (cardItemCount < 12) ? cardItemCount : 12;
                        }

                        let numRows = cardItemsToInclude/4;
                        if (numRows%4 != 0) numRows += 1;
                        const elements = [];
                        let startIndex = 0;
                        let endIndex = 4;
                        let metrics = cardItems;

                        for (let i=0; i<numRows; i++) {
                            let rowMetrics = metrics.splice(startIndex, endIndex);
                            elements.push(
                                <div className="disease-metrics-container-row">
                                    {rowMetrics.map((item, index) => (
                                        Object.entries(item).map(([key, value]) => (
                                            Object.entries(value).map(([keyInner, valueInner]) => {
                                                return  <div className="disease-metrics-card" onClick={() => onShowAlert(value, key)}>
                                                    <p key={keyInner}>
                                                        <strong>{key}</strong><br />
                                                        <strong>Value:</strong> {valueInner.split('-')[0]}<br />
                                                        <strong>Status:</strong> {valueInner.split('-')[1]}<br />
                                                        <strong>Date:</strong> {convertEpochDayToDate(keyInner)}
                                                    </p>
                                                </div>
                                            })
                                        ))
                                    ))}
                                </div>
                            );
                        }
                        return elements;
                    })()
                }
            </div>

            {
                (() => {
                    if (!props.showAllMetrics) {
                        return <Link to="/metrics" className="button">View More</Link>
                     }
                })()
            }
        </div> :
        <div className='track-progress'>
            {
                (!props.showAllMetrics) ?
                <div>
                    <h1>Track Your Progress</h1>
                    <p>This is empty. Login to know more details</p>
                </div> :
                <div>
                    <h1>Metrics</h1>
                    <p>This is empty. Login to know more details</p>
                </div>
            }
        </div>
    );
}

export default TrackProgress;