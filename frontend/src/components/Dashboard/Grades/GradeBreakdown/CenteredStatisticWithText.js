import { Typography } from 'antd';
import React from 'react';
import * as PropTypes from 'prop-types';

export function CenteredStatisticWithText(props) {
  return (
    <div>
      <div align="center">
        <Typography.Title level={1}>{props.stat}</Typography.Title>
      </div>
      <Typography.Text>{props.text}</Typography.Text>
    </div>
  );
}

CenteredStatisticWithText.propTypes = {
  stat: PropTypes.any.isRequired,
  text: PropTypes.string
};
