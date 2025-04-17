import React, { useState, useEffect } from 'react';
    
function Login({ onLogin }) {
  const [username, setUsername] = useState('');

  const handleUsernameChange = (event) => {
    setUsername(event.target.value);
  };

  const saveUsername = () => {
    if (onLogin) {
      onLogin(username); // Pass the new username to the parent component
    }
  };


  return (
    <div>
      <label htmlFor="username">Username: </label>
      <input
        type="text"
        id="username"
        value={username}
        className="search-box"
        placeholder="Enter your username"
        style={{width: "50%",}}
        onChange={handleUsernameChange}
      />
      <button 
      onClick={saveUsername}
      className="user-button"
            style={{
              marginLeft: "10px",
              cursor: "pointer",
            }}
      >Save Username
      </button>
      <p>Saved Username: {username}</p>
    </div>
  );
}

export default Login;