import React from 'react';
import { Link } from 'react-router-dom';
import s from './Header.module.css';

export default function Header({ children }) {
  return (
    <header className={s.header}>
      <Link to="/" className={s.brand}>
        Cointly <span className={s.underline} />
      </Link>
      <div className={s.content}>
        {children}
      </div>
    </header>
  );
}
