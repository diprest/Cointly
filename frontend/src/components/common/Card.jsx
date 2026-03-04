import React from 'react';
export default function Card({ title, icon, onClick, children }){
  return (
    <section className="card">
      {title && <h3 className="card-title">{title}</h3>}
      {icon && <div className="icon-wrap">{icon}</div>}
      {children}
      <button className="btn" onClick={onClick ?? (()=>{})}>
        Открыть <span className="arrow">→</span>
      </button>
    </section>
  );
}
