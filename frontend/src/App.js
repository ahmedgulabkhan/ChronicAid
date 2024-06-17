import { useEffect } from 'react';
import { Routes, Route } from 'react-router-dom';
import { useDispatch } from "react-redux";
import { login } from "./store";

import './App.css';
import Header from './components/header/header-component';
import HomePage from './pages/homepage/homepage-component';
import DiseasesPage from './pages/diseases-page/diseases-page-component';
import CheckSymptomsPage from './pages/check-symptoms-page/check-symptoms-page-component';
import SignInAndSignUpPage from './pages/sign-in-and-sign-up-page/sign-in-and-sign-up-page-component';
import DiseaseDetailsPage from './pages/disease-details-page/disease-details-page-component';
import DiseaseMetricsPage from './pages/disease-metrics-page/disease-metrics-page-component';
import MetricsPage from './pages/metrics-page/metrics-page-component';

function App() {
	const dispatch = useDispatch();
	let status;

	useEffect(() => {
		fetch('http://localhost:8080/api/auth/verify', {
            method: 'GET',
            headers: {
                'Accept': 'application/json',
                'Content-Type': 'application/json'
            },
            credentials: 'include',
        })
        .then(res => {
			status = res.status;
			return res.json();
		}).then(json => {
			if (status === 200) {
				dispatch(login({ username: json.user_name, first_name: json.first_name, isLoggedIn: true }));
			}
		})
        .catch(err => console.log(err));
	}, []);

	return (
    	<div className="App">
			<Header />
			<Routes>
				<Route path='/' element={<HomePage />} />
				<Route path='/diseases' element={<DiseasesPage/>} />
				<Route path='/diseases/:diseaseName' element={<DiseaseDetailsPage/>} />
				<Route path='/diseases/:diseaseName/:metricName' element={<DiseaseMetricsPage/>} />
				<Route path='/metrics' element={<MetricsPage/>} />
				<Route path='/checkSymptoms' element={<CheckSymptomsPage/>} />
				<Route path='/signin' element={<SignInAndSignUpPage/>}/>
			</Routes>
		</div>
  	);
}

export default App;
