import React from "react";
import { Link } from 'react-router-dom';
import { useDispatch, useSelector } from "react-redux";

import './header-styles.scss';
import { logout } from "../../store";

const Header = () => {
    const dispatch = useDispatch();
    const isLoggedIn = useSelector(state => state.user.isLoggedIn);

    const handleLogout = () => {
        fetch('http://localhost:8080/api/auth/logout', {
            method: 'POST',
            headers: {
                'Accept': 'application/json',
                'Content-Type': 'application/json'
            },
            credentials: 'include'
        })
        .then(res => {
            alert("You have logged out of your account!");
            window.location.replace("http://localhost:3000");
        })
        .catch(err => console.log(err));
        dispatch(logout());
    }

    return (
        <div className='header'>
            <div className='container logo-container'>
                <Link to='/' className='logo'>ChronicAid</Link>
            </div>
            <div className='container options'>
                <Link to='/' className='option'>Home</Link>
                <Link to='/diseases' className='option'>Diseases</Link>
                <Link to='/metrics' className='option'>Metrics</Link>
                <Link to='/checkSymptoms' className='option'>Check Symptoms</Link>
            </div>
            <div className='container sign-in-container'>
                {
                    (isLoggedIn === true) ? 
                    (<Link className='log-out' onClick={handleLogout}>Logout</Link>)
                    :
                    (<Link to='/signin' className='sign-in'>Sign In / Register</Link>)
                }
            </div>
        </div>
    );
}

export default Header;