import React from 'react';
import { useNavigate } from 'react-router-dom';
import { useSelector } from 'react-redux';
import Card from '../../components/common/Card';
import ChartIcon from '../../components/icons/ChartIcon';
import UserIcon from '../../components/icons/UserIcon';

export default function Home() {
  const navigate = useNavigate();
  const isAuth = useSelector((state) => state.user.isAuth);

  const handleNavigation = (path) => {
    if (isAuth) {
      navigate(path);
    } else {
      navigate('/auth');
    }
  };

  return (
    <div className="page">
      <div className="center">
        <div className="hero">
          <h1 className="logo">
            <span className="logo-under">Coint</span>ly
          </h1>
          <p className="tagline">Лёгкий способ быть в плюсе</p>
        </div>

        <div className="cards">
          <Card
            title="Все монеты"
            icon={<ChartIcon />}
            onClick={() => handleNavigation('/coins')}
          />
          <Card
            title="Профиль"
            icon={<UserIcon />}
            onClick={() => handleNavigation('/profile')}
          />
        </div>
      </div>
    </div>
  );
}
