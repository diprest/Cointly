import React, { useState } from 'react';
import { Link, useLocation } from 'react-router-dom';
import s from './ProfileNavigation.module.css';
import SettingsModal from '../SettingsModal/SettingsModal';

export default function ProfileNavigation() {
  const location = useLocation();
  const path = location.pathname;
  const [isSettingsOpen, setIsSettingsOpen] = useState(false);

  const getBtnClass = (targetPath) => {
    return path === targetPath ? s.btnPrimary : s.btnGhost;
  };

  return (
    <>
      <nav className={s.nav}>
        <Link to="/coins" className={s.btnGhost}>Маркет</Link>
        <Link to="/profile" className={getBtnClass('/profile')}>Обзор</Link>
        <Link to="/profile-wallet" className={getBtnClass('/profile-wallet')}>Мой кошелёк</Link>
        <Link to="/profile-history" className={getBtnClass('/profile-history')}>История</Link>
        <button className={s.btnGhost} onClick={() => setIsSettingsOpen(true)}>Настройки</button>
      </nav>
      <SettingsModal isOpen={isSettingsOpen} onClose={() => setIsSettingsOpen(false)} />
    </>
  );
}
