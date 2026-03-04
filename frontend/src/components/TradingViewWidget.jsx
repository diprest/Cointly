import React, { useEffect, useRef } from 'react';

export default function TradingViewWidget({ symbol }) {
  const containerRef = useRef(null);

  useEffect(() => {
    const scriptId = 'tradingview-widget-script';
    const containerId = `tradingview_widget_${Math.random().toString(36).substring(7)}`;

    if (containerRef.current) {
      containerRef.current.id = containerId;
    }

    const initWidget = () => {
      if (window.TradingView && containerRef.current) {
        new window.TradingView.widget({
          autosize: true,
          symbol: `BINANCE:${symbol}`,
          interval: "D",
          timezone: "Etc/UTC",
          theme: "dark",
          style: "1",
          locale: "ru",
          toolbar_bg: "#f1f3f6",
          enable_publishing: false,
          allow_symbol_change: true,
          container_id: containerId,
        });
      }
    };

    if (!document.getElementById(scriptId)) {
      const script = document.createElement('script');
      script.id = scriptId;
      script.src = 'https://s3.tradingview.com/tv.js';
      script.type = 'text/javascript';
      script.onload = initWidget;
      document.head.appendChild(script);
    } else {
      initWidget();
    }
  }, [symbol]);

  return (
    <div className='tradingview-widget-container' style={{ height: "100%", width: "100%" }}>
      <div ref={containerRef} style={{ height: "100%", width: "100%" }} />
    </div>
  );
}
