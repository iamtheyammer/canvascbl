import React from 'react';
import PropTypes from 'prop-types';
import { Spin, Typography } from 'antd';

function Loading(props) {
  const { text } = props;

  return (
    <div align="center">
      <Spin />
      <span style={{ paddingTop: '20px' }} />
      <Typography.Title level={3}>{`Loading ${text}...`}</Typography.Title>
    </div>
  );
}

Loading.propTypes = {
  text: PropTypes.string
};

export default Loading;
