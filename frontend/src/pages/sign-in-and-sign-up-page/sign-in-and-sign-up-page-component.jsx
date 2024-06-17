import React from 'react';
import { useState } from 'react';
import { useDispatch } from "react-redux";
import { login } from '../../store';
import { useNavigate } from 'react-router-dom';

import './sign-in-and-sign-up-page-styles.scss';

const SignInAndSignUpPage = () => {
    const navigate = useNavigate();
    const dispatch = useDispatch();
    const [accountStatus, setAccountStatus] = useState("sign-in");
    const [signInInputValues, setSignInInputValues] = useState({
        email: '',
        password: ''
    });
    const [signUpInputValues, setSignUpInputValues] = useState({
        firstName: '',
        lastName: '',
        userName: '',
        email: '',
        password: ''
    });
    const [errorMessageSignIn, setErrorMessageSignIn] = useState("");
    const [errorMessageSignUp, setErrorMessageSignUp] = useState("");

    const handleSignInInputChange = (e) => {
        const { name, value } = e.target;
        setSignInInputValues({
            ...signInInputValues,
            [name]: value
        });
    }
    
    const handleSignUpInputChange = (e) => {
        const { name, value } = e.target;
        setSignUpInputValues({
            ...signUpInputValues,
            [name]: value
        });
    }
    
    const handleSignInSubmit = (e) => {
        e.preventDefault();
        const data = {
            email: signInInputValues.email,
            password: signInInputValues.password
        }

        fetch('http://localhost:8080/api/auth/login', {
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
            if (json.error) {
                setErrorMessageSignIn(json.error);
            } else {
                dispatch(login({ username: json.user_name, first_name: json.first_name, isLoggedIn: true }));
                navigate("/");
            }
        })
        .catch(err => setErrorMessageSignIn(err));
    }

    const handleSignUpSubmit = (e) => {
        e.preventDefault();
        const data = {
            first_name: signUpInputValues.firstName,
            last_name: signUpInputValues.lastName,
            user_name: signUpInputValues.userName,
            email: signUpInputValues.email,
            password: signUpInputValues.password
        }
        fetch('http://localhost:8080/api/auth/signup', {
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
            if (json.error) {
                setErrorMessageSignUp(json.error);
            } else {
                dispatch(login({ username: json.user_name, first_name: json.first_name, isLoggedIn: true }))
                navigate("/");
            }
        })
        .catch(err => setErrorMessageSignUp(err));
    }

    const handleToggleAccountStatus = () => {
        if (accountStatus === "sign-in") {
            setAccountStatus("sign-up");
        } else {
            setAccountStatus("sign-in");
        }
    }

    return (
        <div className="sign-in-sign-up-page">
            {
                (accountStatus === "sign-in") ?
                <form onSubmit={handleSignInSubmit}>
                    <h1>Sign In</h1>
                    <div className="input-field">
                        <input
                        type="text"
                        name="email"
                        placeholder="Email"
                        value={signInInputValues.email}
                        onChange={handleSignInInputChange}
                        />
                    </div>
                    <div className="input-field">
                        <input
                        type="text"
                        name="password"
                        placeholder="Password"
                        value={signInInputValues.password}
                        onChange={handleSignInInputChange}
                        />
                    </div>
                    <button type="submit">Submit</button>
                    {
                        (() => {
                            if (errorMessageSignIn !== "") {
                                return <p className="error-message">{errorMessageSignIn}</p>
                            }
                        })()
                    }
                    <p onClick={handleToggleAccountStatus}>Don't have an account? Sign Up</p>
                </form> :
                <form onSubmit={handleSignUpSubmit}>
                    <h1>Sign Up</h1>
                    <div className="input-field">
                        <input
                        type="text"
                        name="firstName"
                        placeholder="First Name"
                        value={signUpInputValues.firstName}
                        onChange={handleSignUpInputChange}
                        />
                    </div>
                    <div className="input-field">
                        <input
                        type="text"
                        name="lastName"
                        placeholder="Last Name"
                        value={signUpInputValues.lastName}
                        onChange={handleSignUpInputChange}
                        />
                    </div>
                    <div className="input-field">
                        <input
                        type="text"
                        name="userName"
                        placeholder="Username"
                        value={signUpInputValues.userName}
                        onChange={handleSignUpInputChange}
                        />
                    </div>
                    <div className="input-field">
                        <input
                        type="text"
                        name="email"
                        placeholder="Email"
                        value={signUpInputValues.email}
                        onChange={handleSignUpInputChange}
                        />
                    </div>
                    <div className="input-field">
                        <input
                        type="text"
                        name="password"
                        placeholder="Password"
                        value={signUpInputValues.password}
                        onChange={handleSignUpInputChange}
                        />
                    </div>
                    <button type="submit">Submit</button>
                    {
                        (() => {
                            if (errorMessageSignUp !== "") {
                                return <p className="error-message">{errorMessageSignUp}</p>
                            }
                        })()
                    }
                    <p onClick={handleToggleAccountStatus}>Already have an account? Sign In</p>
                </form>
            }

        </div>
    );
}

export default SignInAndSignUpPage;