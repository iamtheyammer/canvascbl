import React from 'react';
import PropTypes from 'prop-types';

function Padding(props) {
  const { top, right, bottom, left, all, br } = props;
  return br ? (
    <br />
  ) : (
    <div style={{ padding: all ? all : `${top} ${right} ${bottom} ${left}` }} />
  );
}

Padding.propTypes = {
  top: PropTypes.oneOfType([PropTypes.string, PropTypes.number]),
  right: PropTypes.oneOfType([PropTypes.string, PropTypes.number]),
  bottom: PropTypes.oneOfType([PropTypes.string, PropTypes.number]),
  left: PropTypes.oneOfType([PropTypes.string, PropTypes.number]),
  all: PropTypes.oneOfType([PropTypes.string, PropTypes.number]),
  br: PropTypes.bool
};

Padding.defaultProps = {
  top: 0,
  right: 0,
  bottom: 0,
  left: 0,
  br: false
};

export default Padding;
