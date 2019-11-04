import React, { useState } from 'react';
import { Typography, Card, Icon, Skeleton } from 'antd';
import { Tabs as MobileTabs } from 'antd-mobile';
import * as PropTypes from 'prop-types';
import { isMobile } from 'react-device-detect';
import { CenteredStatisticWithText } from './CenteredStatisticWithText';
import { gradeMapByGrade } from '../../../../util/canvas/calculateGradeFromOutcomes';
import { ReactComponent as plusIcon } from '../../../../assets/plus.svg';

const tabList = [
  {
    key: 'userGrade',
    tab: 'Your Grade'
  },
  {
    key: 'averageGrade',
    tab: (
      <div>
        <Icon component={plusIcon} />
        {isMobile && ' '}
        <Typography.Text>Average Grade</Typography.Text>
      </div>
    )
  }
];

const mobileTabList = tabList.map(t => ({ title: t.tab, sub: t.key }));

function GradeCard(props) {
  const [activeTabKey, setActiveTabKey] = useState(tabList[0].key);

  const { currentGrade, averageGrade, userHasValidSubscription } = props;

  function generateCardContent(key) {
    switch (key) {
      case 'userGrade':
        const { min, max } = gradeMapByGrade[currentGrade];
        return (
          <CenteredStatisticWithText
            stat={currentGrade}
            text={`Your current grade, ${currentGrade}, requires 75% of outcomes to be
              above ${max} and no outcomes to be below ${min}.`}
          />
        );
      case 'averageGrade':
        if (!userHasValidSubscription) {
          return (
            <div>
              <Typography.Title level={3}>CanvasCBL+ Required</Typography.Title>
              <Typography.Text>
                You need CanvasCBL+ to use this feature. Click on the 'Upgrades'
                page to check it out and upgrade!
              </Typography.Text>
            </div>
          );
        }

        if (!averageGrade) {
          return <Skeleton active />;
        }

        const { numFactors, averageGrade: grade, error } = averageGrade;

        if (numFactors === -1 && error === 'not enough factors') {
          return (
            <CenteredStatisticWithText
              stat={'N/A'}
              text={
                'An average grade for this class is not available because there are not enough students ' +
                'in this class who have logged into CanvasCBL in the last 24 hours. Encourage your classmates to log in!'
              }
            />
          );
        }

        if (error) {
          return (
            <CenteredStatisticWithText
              stat={'N/A'}
              text={'There was an error retrieving the average grade.'}
            />
          );
        }

        return (
          <CenteredStatisticWithText
            stat={grade}
            text={`The average grade in this class is ${grade}, which was calculated from ${numFactors} data points.`}
          />
        );
      default:
        return null;
    }
  }

  if (isMobile) {
    return (
      <MobileTabs tabs={mobileTabList} initialPage={activeTabKey}>
        {mobileTabList.map(t => (
          <div key={t.sub}>{generateCardContent(t.sub)}</div>
        ))}
      </MobileTabs>
    );
  }

  return (
    <Card
      tabList={tabList}
      activeTabKey={activeTabKey}
      onTabChange={setActiveTabKey}
    >
      {generateCardContent(activeTabKey)}
    </Card>
  );
}

GradeCard.propTypes = {
  currentGrade: PropTypes.string.isRequired,
  userHasValidSubscription: PropTypes.bool.isRequired,
  averageGrade: PropTypes.shape({
    numFactors: PropTypes.number.isRequired,
    averageGrade: PropTypes.string.isRequired,
    error: PropTypes.string
  })
};

export default GradeCard;
