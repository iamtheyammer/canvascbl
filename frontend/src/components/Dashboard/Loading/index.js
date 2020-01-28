import React from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';
import { Spin, Typography } from 'antd';
import useInterval from '@use-hooks/interval';
import { updateNumberOfDots } from '../../../actions/components/loading';

function Loading(props) {
  const { dots, text, dispatch } = props;

  useInterval(() => {
    let newNumDots;
    if (dots === 3) {
      newNumDots = 0;
    } else {
      newNumDots = dots !== undefined && dots !== null ? dots + 1 : 1;
    }
    dispatch(updateNumberOfDots(newNumDots));
  }, 500);

  let dotStr = '';
  for (let i = 0; i < dots; i++) {
    dotStr += '.';
  }

  return (
    <div align="center">
      <Spin />
      <span style={{ paddingTop: '20px' }} />
      <Typography.Title
        level={3}
      >{`Loading ${text}${dotStr}`}</Typography.Title>
    </div>
  );
}

Loading.propTypes = {
  text: PropTypes.string
};

export default connect(state => ({
  dots: state.components.loading.dots
}))(Loading);
