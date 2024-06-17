import React, { useState } from "react";
import { Link } from "react-router-dom";
import { useEffect } from "react";
import { useSelector } from "react-redux";

import './your-diseases-styles.scss';

const YourDiseases = (props) => {
    const username = useSelector(state => state.user.username);
    const isLoggedIn = useSelector(state => state.user.isLoggedIn);
    const [userDiseases, setUserDiseases] = useState([]);
    const Swal = require('sweetalert2');

    useEffect(() => {
        if (username === undefined || username === "") return;
        fetch('http://localhost:8080/api/diseases/' + username, {
            method: 'GET',
            headers: {
                'Accept': 'application/json',
                'Content-Type': 'application/json'
            },
            credentials: 'include',
        })
        .then(res => res.json())
        .then(json => setUserDiseases(json))
        .catch(err => console.log(err));
    }, [username]);

    const toUpperCamelCase = (str) => {
        return str
            .split('-')
            .map(word => word.charAt(0).toUpperCase() + word.slice(1).toLowerCase())
            .join(' ');
    }

    const handleAddNewDisease = () => {
        Swal.fire({
            title: "Add New Disease",
            input: "text",
            inputPlaceholder: "Disease Name"
        }).then((result) => {
            if (result.isConfirmed && result.value !== "") {
                const data = {
                    disease_name: result.value
                }
                fetch('http://localhost:8080/api/diseases/' + username, {
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
                    console.log(json);
                    window.location.reload();
                })
                .catch(err => console.log(err));
            }
        });
    }

    return (
        (isLoggedIn === true) ?
        <div className='your-diseases' >
            <div className="heading-section">
                <h1>Your Diseases</h1>
                <button className="add-new-button" onClick={handleAddNewDisease}>Add New Disease</button>
            </div>
            <div className="diseases-container-wrapper">
                {
                    (() => {
                        let diseases = (props.showAllDiseases) ? userDiseases.slice(0, 6) : userDiseases;
                        let numRows = diseases.length/3;
                        if (numRows%3 != 0) numRows += 1;
                        const elements = [];
                        let startIndex = 0;
                        let endIndex = 3;

                        for (let i=0; i<numRows; i++) {
                            let rowDiseases = diseases.splice(startIndex, endIndex);
                            elements.push(
                                <div className="diseases-container-row">
                                    {rowDiseases.map(item => (
                                        <Link to={`/diseases/${item.toLowerCase()}`} className="diseases-card-link">
                                            <div className="diseases-card">
                                                <strong>{toUpperCamelCase(item)}</strong>
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

            {
                (() => {
                    if (!props.showAllDiseases) {
                        return <Link to="/diseases" className="button">View More</Link>
                    }
                })()
            }

        </div> :
        <div className='your-diseases'>
            <h1>Your Diseases</h1>
            <p>This is empty. Login to know more details</p>
        </div>
    );
}

export default YourDiseases;