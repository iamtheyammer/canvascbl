import React, { useState } from 'react';
import { Card, Icon, Typography } from 'antd';
import { Tabs as MobileTabs } from 'antd-mobile';
import { isMobile } from 'react-device-detect';
import calculateMeanAverage from '../../../../util/calculateMeanAverage';
import roundNumberToDigits from '../../../../util/roundNumberToDigits';
import v4 from 'uuid/v4';
import { gradeMapByGrade } from '../../../../util/canvas/gradeMapByGrade';
import { CenteredStatisticWithText } from './CenteredStatisticWithText';
import { ReactComponent as plusIcon } from '../../../../assets/plus.svg';
import { tabImplementations, trackTabChange } from '../../../../util/tracking';

const tabList = [
  {
    key: 'lowestOutcome',
    tab: 'Lowest Outcome'
  },
  {
    key: 'averageOutcomeScore',
    tab: 'Average Outcome Score'
  },
  {
    key: 'toGetAnA',
    tab: (
      <div>
        <Icon component={plusIcon} />
        {isMobile && ' '}
        <Typography.Text>How To Get An A</Typography.Text>
      </div>
    )
  },
  {
    key: 'moreInfo',
    tab: 'More Info'
  }
];

const mobileTabList = tabList.map(t => ({ title: t.tab, sub: t.key }));

function OutcomeInfo(props) {
  const [activeTabKey, setActiveTabKey] = useState(tabList[0].key);

  const {
    lowestOutcome,
    outcomeRollupScores,
    grade,
    userHasValidSubscription
  } = props;

  const { min: AMin, max: AMax } = gradeMapByGrade['A'];
  const outcomeScores = Object.values(outcomeRollupScores)
    .map(or => or.average)
    .sort((a, b) => a - b);
  const seventyFivePercentOfOutcomes = Math.floor(
    (75 * outcomeScores.length) / 100
  );

  function generateCardContent(key) {
    switch (key) {
      case 'lowestOutcome':
        return (
          <CenteredStatisticWithText
            stat={+lowestOutcome.rollupScore.average.toFixed(2)}
            text={`Your lowest outcome is ${lowestOutcome.outcome
              .display_name ||
              lowestOutcome.outcome
                .title}, with a score of ${+lowestOutcome.rollupScore.average.toFixed(
              2
            )}.`}
          />
        );
      case 'averageOutcomeScore':
        const meanOutcomeScore = roundNumberToDigits(
          calculateMeanAverage(
            Object.values(outcomeRollupScores).map(or => or.average)
          ),
          3
        );
        return (
          <CenteredStatisticWithText
            stat={meanOutcomeScore}
            text={`Your average outcome score is ${meanOutcomeScore}.`}
          />
        );
      case 'toGetAnA':
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

        if (grade.grade.grade === 'A') {
          return (
            <div>
              <div align="center">
                <Typography.Title level={1}>Nice job!</Typography.Title>
              </div>
              <Typography.Text>
                You've already got an A in this class! To keep it, make sure 75%
                of outcome scores are over {AMax} and that no outcome scores
                drop below {AMin}.
              </Typography.Text>
            </div>
          );
        } else {
          return (
            <div>
              <Typography.Text>
                To get an A in this class, you'll need to make sure that:
              </Typography.Text>
              <ul>
                <li>
                  <Typography.Text
                    delete={lowestOutcome.lowestCountedOutcome > AMax}
                  >
                    {seventyFivePercentOfOutcomes}/{outcomeScores.length}{' '}
                    outcomes are above {AMax} (currently,{' '}
                    {outcomeScores.filter(o => o < AMax).length} outcomes are
                    not above {AMax})
                  </Typography.Text>
                </li>
                <li>
                  <Typography.Text delete={outcomeScores[0] >= AMin}>
                    No outcomes are below {AMin} (currently,{' '}
                    {outcomeScores.filter(o => o < AMin).length} outcomes are
                    below {AMin})
                  </Typography.Text>
                </li>
              </ul>
              Please note that this section may not be technically accurate as
              your lowest outcome score is dropped.
            </div>
          );
        }
      case 'moreInfo':
        const cardWithContent = content => (
          <Card.Grid key={v4()} hoverable="false">
            {content}
          </Card.Grid>
        );

        return [
          cardWithContent(
            `75% (rounded) of ${outcomeScores.length} (number of outcomes with a grade) is ${seventyFivePercentOfOutcomes}.`
          ),
          cardWithContent(`More info is coming to this section in the future.`)
        ];
      default:
        return (
          <Typography.Text>
            There was an error: OutcomeInfo Default Case Used
          </Typography.Text>
        );
    }
  }

  if (isMobile) {
    return (
      <MobileTabs
        tabs={mobileTabList}
        initialPage={activeTabKey}
        className="mobile-tabs"
        renderTabBar={props => (
          <MobileTabs.DefaultTabBar {...props} page={1.65} />
        )}
        onChange={data => {
          trackTabChange(
            tabImplementations.outcomeInfo.name,
            tabImplementations.outcomeInfo.tabNames[data.sub]
          );
        }}
      >
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
      onTabChange={newTabKey => {
        setActiveTabKey(newTabKey);
        trackTabChange(
          tabImplementations.outcomeInfo.name,
          tabImplementations.outcomeInfo.tabNames[newTabKey]
        );
      }}
    >
      {generateCardContent(activeTabKey)}
    </Card>
  );
}

export default OutcomeInfo;
