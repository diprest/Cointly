import React from 'react';
import { Routes, Route, Navigate, useLocation } from 'react-router-dom';
import { useSelector } from 'react-redux';
import './App.css';

import Home from '../pages/Home/Home';
import Auth from '../pages/Auth/Auth';
import Coins from '../pages/Coins/Coins';
import Profile from '../pages/Profile/Profile';
import ProfileWallet from '../pages/ProfileWallet/ProfileWallet';
import ProfileHistory from '../pages/ProfileHistory/ProfileHistory';

function RequireAuth({ children }) {
  const isAuth = useSelector((state) => state.user.isAuth);
  const location = useLocation();

  if (!isAuth) {
    return <Navigate to="/auth" state={{ from: location }} replace />;
  }

  return children;
}

export default function App() {
  return (
    <Routes>
      <Route path="/" element={<Home />} />
      <Route path="/auth" element={<Auth />} />

      <Route
        path="/coins"
        element={
          <RequireAuth>
            <Coins />
          </RequireAuth>
        }
      />
      <Route
        path="/profile"
        element={
          <RequireAuth>
            <Profile />
          </RequireAuth>
        }
      />
      <Route
        path="/wallet"
        element={
          <RequireAuth>
            <ProfileWallet />
          </RequireAuth>
        }
      />
      <Route
        path="/history"
        element={
          <RequireAuth>
            <ProfileHistory />
          </RequireAuth>
        }
      />
    </Routes>
  );
}
