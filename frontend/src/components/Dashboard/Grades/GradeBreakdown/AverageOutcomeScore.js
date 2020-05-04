import React, { useState } from 'react';
import { connect } from 'react-redux';
import { Skeleton, Typography } from 'antd';
import v4 from 'uuid/v4';
import { getAverageScoreForOutcome } from '../../../../actions/plus';

function AverageOutcomeScore(props) {
  const { dispatch, plus, outcomeId, loading, error } = props;

  const [getAverageOutcomeScoreId, setGetAverageOutcomeScoreId] = useState('');
  const averageIsLoading = loading.includes(getAverageOutcomeScoreId);
  const getAverageError = error[getAverageOutcomeScoreId];

  if (!plus.session.has_valid_subscription) {
    return (
      <Typography.Text>
        You need CanvasCBL+ to use this feature. Go to the Upgrades page at the
        top to check it out and upgrade!
      </Typography.Text>
    );
  }

  if (
    !getAverageOutcomeScoreId &&
    (!plus.outcomeAverages || !plus.outcomeAverages[outcomeId])
  ) {
    const id = v4();
    dispatch(getAverageScoreForOutcome(id, outcomeId));
    setGetAverageOutcomeScoreId(id);
    return null;
  }

  if (getAverageError) {
    return (
      <Typography.Text>
        There was an error getting the average score for this outcome.
      </Typography.Text>
    );
  }

  if (
    averageIsLoading ||
    !plus.outcomeAverages ||
    !plus.outcomeAverages[outcomeId]
  ) {
    return (
      <Skeleton lines={1} title={{ width: '50%' }} paragraph={false} active />
    );
  }

  const avg = plus.outcomeAverages[outcomeId];

  if (avg.error === 'not enough factors') {
    return (
      <Typography.Text>
        Not enough factors to calculate. Encourage classmates to sign in to
        CanvasCBL!
      </Typography.Text>
    );
  }

  if (avg.error) {
    return (
      <Typography.Text>
        There was an error getting the average: {avg.error}
      </Typography.Text>
    );
  }

  return (
    <Typography.Text>
      The average score for this outcome is{' '}
      <Typography.Text strong>{avg.averageScore}</Typography.Text>, calculated
      from {avg.numFactors} data points.
    </Typography.Text>
  );
}

const ConnectedAverageOutcomeScore = connect((state) => ({
  loading: state.loading,
  error: state.error,
  plus: state.plus
}))(AverageOutcomeScore);

export default ConnectedAverageOutcomeScore;
