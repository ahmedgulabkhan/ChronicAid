import React, { useState } from "react";

import "./check-symptoms-page-styles.scss";

const ActivityPage = () => {
    const [text, setText] = useState('');
    const Swal = require('sweetalert2');

    const handleTextChange = (e) => {
        setText(e.target.value);
    };

    const handleSubmit = (e) => {
        e.preventDefault();

        const data = {
            message: text
        }
        
        fetch('http://localhost:8080/api/diseases/symptoms', {
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
    };

    return (
        <div className="check-symptoms-page">
            <h1>Check Symptoms Using AI</h1>
            <div className="text-area-section">
                <form onSubmit={handleSubmit}>
                    <div className="input-field">
                        <textarea
                            name="text"
                            placeholder="Enter your symptoms here..."
                            value={text}
                            onChange={handleTextChange}
                        />
                    </div>
                    <button type="submit">Submit</button>
                </form>
            </div>
        </div>
    );
}

export default ActivityPage;
