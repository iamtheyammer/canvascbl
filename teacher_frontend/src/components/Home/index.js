import React from 'react';
import { Redirect } from 'react-router-dom';

function Home(props) {
  return <Redirect to="/dashboard" />;
}

export default Home;
