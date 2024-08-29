import React from 'react';
import { useReactFlow } from 'reactflow';

const ZoomButtons = ({ size = 40, color = '#ffffff', backgroundColor = '#808080' }) => {
  const { zoomIn, zoomOut } = useReactFlow();

  const buttonStyle = {
    width: `${size}px`,
    height: `${size}px`,
    borderRadius: '50%',
    backgroundColor,
    border: 'none',
    cursor: 'pointer',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
    margin: '5px',
    padding: 0,
  };

  const containerStyle = {
    position: 'absolute',
    right: '10px',
    bottom: '10px',
    display: 'flex',
    flexDirection: 'column',
    zIndex: 4,
  };

  return (
    <div style={containerStyle}>
      <button style={buttonStyle} onClick={() => zoomIn()}>
        <svg width={size/2} height={size/2} viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
          <path d="M12 4V20M4 12H20" stroke={color} strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"/>
        </svg>
      </button>
      <button style={buttonStyle} onClick={() => zoomOut()}>
        <svg width={size/2} height={size/2} viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
          <path d="M4 12H20" stroke={color} strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"/>
        </svg>
      </button>
    </div>
  );
};

export default ZoomButtons;
