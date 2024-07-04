'use client';

import React from 'react';
import SignupForm from '../pages/users/signup';
import LoginForm from '../pages/users/login';
import { BrowserRouter, Route, Routes } from "react-router-dom"
import {Row} from "antd";

const Signup = () => (
    <div>
        <Row justify="center" align="middle" style={{ display: 'flex' }}>
            <SignupForm />
        </Row>
    </div>
);


const Login = () => (
    <div>
        <Row justify="center" align="middle" style={{ display: 'flex' }}>
            <LoginForm />
        </Row>
    </div>
);

const App = () => {
    return (
        <BrowserRouter>
            <Routes>
                <Route index element={<Signup />} />
                <Route path="signup" element={<Signup />} />
                <Route path="login" element={<Login />} />
            </Routes>
        </BrowserRouter>
    );
}

export default App;