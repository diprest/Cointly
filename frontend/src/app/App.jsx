import React from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate, useLocation } from 'react-router-dom';
import { Provider, useSelector } from 'react-redux';
import { ToastContainer } from 'react-toastify';
import 'react-toastify/dist/ReactToastify.css';
import { store } from '../store/store';
import Home from '../pages/Home/Home';
import Auth from '../pages/Auth/Auth';
import Coins from '../pages/Coins/Coins';
import Coin from '../pages/Coin/Coin';
import Profile from '../pages/Profile/Profile';
import ProfileHistory from '../pages/ProfileHistory/ProfileHistory';
import ProfileWallet from '../pages/ProfileWallet/ProfileWallet';
import '../App.css';

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
    <Provider store={store}>
      <Router>
        <ToastContainer position="top-right" theme="dark" />
        <Routes>
          <Route path="/" element={<Home />} />
          <Route path="/auth" element={<Auth />} />

          {/* Protected Routes */}
          <Route
            path="/coins"
            element={
              <RequireAuth>
                <Coins />
              </RequireAuth>
            }
          />
          <Route
            path="/coin/:symbol"
            element={
              <RequireAuth>
                <Coin />
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
            path="/profile-history"
            element={
              <RequireAuth>
                <ProfileHistory />
              </RequireAuth>
            }
          />
          <Route
            path="/profile-wallet"
            element={
              <RequireAuth>
                <ProfileWallet />
              </RequireAuth>
            }
          />
          <Route path="/wallet" element={<Navigate to="/profile-wallet" replace />} />
          <Route path="/history" element={<Navigate to="/profile-history" replace />} />
        </Routes>
      </Router>
    </Provider>
  );
}
