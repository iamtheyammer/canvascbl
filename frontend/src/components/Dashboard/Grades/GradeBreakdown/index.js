import React, { useState, useEffect } from 'react';
import { connect } from 'react-redux';
import { Redirect, Link } from 'react-router-dom';
import v4 from 'uuid/v4';
import moment from 'moment';
import { isMobile } from 'react-device-detect';
import { flatten } from 'lodash';

import {
  Typography,
  notification,
  Row,
  Col,
  Button,
  Table,
  Icon,
  Tag
} from 'antd';

import {
  WhiteSpace as MobileWhiteSpace,
  Accordion as MobileAccordion,
  List as MobileList
} from 'antd-mobile';

import './index.css';

import {
  getAssignmentsForCourse,
  getOutcomeAlignmentsForCourse,
  changeActiveUser,
  getIndividualOutcome
} from '../../../../actions/canvas';
import { dateAsc } from '../../../../util/sort';

import { desc } from '../../../../util/sort';
import ConnectedErrorModal from '../../ErrorModal';
import { ReactComponent as PopOutIcon } from '../../../../assets/pop_out.svg';
import { ReactComponent as PlusIcon } from '../../../../assets/plus.svg';

import OutcomeInfo from './OutcomeInfo';
import PopoutLink from '../../../PopoutLink';
import GradeCard from './GradeCard';
import { getAverageGradeForCourse } from '../../../../actions/plus';
import ConnectedAverageOutcomeScore from './AverageOutcomeScore';
import FutureAssignmentsForOutcome from './FutureAssignmentsForOutcome';
import Loading from '../../Loading';

const outcomeTableColumns = [
  {
    title: 'Name',
    dataIndex: 'name',
    key: 'name',
    sorter: (a, b) => desc(a.name, b.name),
    render: text => (
      <div>
        <Typography.Text>{text}</Typography.Text>
        <span style={{ width: '7px', display: 'inline-block' }} />
      </div>
    )
  },
  {
    title: 'Score',
    dataIndex: 'score',
    key: 'score',
    sorter: (a, b) => a.score - b.score
  },
  {
    title: 'Lowest Score Dropped',
    dataIndex: 'worstScoreDropped',
    key: 'worstScoreDropped',
    render: didDrop => (
      <Tag color={didDrop ? 'green' : 'red'}>{didDrop ? 'Yes' : 'No'}</Tag>
    ),
    sorter: (a, b) => {
      const A = a.worstScoreDropped;
      const B = b.worstScoreDropped;
      if (A === B) {
        return 0;
      } else if (A && !B) {
        return 1;
      }

      return -1;
    }
  },
  {
    title: 'Times Assessed',
    dataIndex: 'timesAssessed',
    key: 'timesAssessed',
    sorter: (a, b) => a.timesAssessed - b.timesAssessed
  },
  {
    title: 'Last Assignment',
    dataIndex: 'lastAssignment',
    key: 'lastAssignment',
    sorter: (a, b) => desc(a.lastAssignment, b.lastAssignment)
  }
];

const assignmentTableOutcomes = [
  {
    title: 'Assignment Name',
    dataIndex: 'assignmentName',
    key: 'assignmentName'
  },
  {
    title: 'Score',
    dataIndex: 'score',
    key: 'score',
    render: score => `${score.score}/${score.possible} (${score.percent}%)`
  },
  {
    title: 'Last Submission',
    dataIndex: 'lastSubmission',
    key: 'lastSubmission'
  },
  {
    title: 'Mastery Reached',
    dataIndex: 'masteryReached',
    key: 'masteryReached',
    render: mastery => (
      <div style={{ margin: 'auto' }}>
        {mastery === true ? <Icon type="check" /> : <Icon type="close" />}
      </div>
    )
  },
  {
    title: 'Actions',
    key: 'actions',
    render: (text, record) => (
      <div>
        <PopoutLink url={record.assignmentUrl}>
          Open on Canvas <Icon component={PopOutIcon} />
        </PopoutLink>
      </div>
    )
  }
];

function GradeBreakdown(props) {
  const [getOutcomesIds, setGetOutcomesIds] = useState([]);
  const [getAssignmentsId, setGetAssignmentsId] = useState('');
  const [getOutcomeAlignmentsId, setGetOutcomeAlignmentsId] = useState('');
  const [getPlusAverageId, setGetPlusAverageId] = useState('');

  const [loadingText, setLoadingText] = useState('');

  const {
    dispatch,
    loading,
    grades,
    error,
    activeUserId,
    users,
    courses,
    outcomes,
    outcomeResults,
    outcomeAlignments,
    observees,
    assignments,
    session,
    gradeAverages
  } = props;

  const err =
    error[getOutcomeAlignmentsId] || getOutcomesIds.filter(id => error[id])[0];

  const courseId = parseInt(props.match.params.courseId);

  const activeUser = users && activeUserId && users[activeUserId];

  const allOutcomes = props.outcomes;

  useEffect(
    () => {
      if (isNaN(courseId)) return;

      // we can display the page without the average loading
      if (
        !getPlusAverageId &&
        session &&
        session.has_valid_subscription &&
        (!gradeAverages || !gradeAverages[courseId])
      ) {
        const id = v4();
        dispatch(getAverageGradeForCourse(id, courseId));
        setGetPlusAverageId(id);
      }

      const course = props.courses.filter(c => c.id === courseId)[0];

      if (!course) {
        return;
      }

      const courseGradedUsers = course.enrollments.map(
        e => e.associated_user_id || e.user_id
      );

      if (!courseGradedUsers.includes(activeUserId)) {
        return;
      }

      // we can display the page without loading alignments
      if (
        !getOutcomeAlignmentsId &&
        activeUserId &&
        (!outcomeAlignments || !outcomeAlignments[courseId])
      ) {
        const id = v4();
        dispatch(getOutcomeAlignmentsForCourse(id, courseId, activeUserId));
        setGetOutcomeAlignmentsId(id);
      }

      if (!grades[activeUserId][courseId].averages) {
        return;
      }

      const neededOutcomes = Object.keys(
        grades[activeUserId][courseId].averages
      );

      // if active user and courses and no outcomes are loading and not every needed outcome had
      if (
        activeUser &&
        courses &&
        !getOutcomesIds.length &&
        (!getOutcomesIds.some(oId => loading.includes(oId)) &&
          !neededOutcomes.every(
            id => allOutcomes && allOutcomes.some(o => o.id === parseInt(id))
          ))
      ) {
        neededOutcomes.forEach(noId => {
          const id = v4();
          dispatch(getIndividualOutcome(id, noId));
          setLoadingText('outcomes');
          setGetOutcomesIds([...getOutcomesIds, id]);
        });
      }

      if (!courseGradedUsers.includes(activeUserId)) {
        return;
      }

      if (
        activeUser &&
        activeUser.id &&
        (!assignments || !assignments[courseId]) &&
        !getAssignmentsId
      ) {
        const id = v4();
        dispatch(getAssignmentsForCourse(id, courseId));
        setLoadingText('your assignments');
        setGetAssignmentsId(id);
      }
    },
    // disabling because we specifically only want to re-run this on a props change
    // eslint-disable-next-line
    [props]
  );

  if (isNaN(courseId)) {
    notification.error({
      message: 'Invalid Course ID',
      description: 'Course IDs contain only numbers.'
    });
    return <Redirect to="/dashboard/grades" />;
  }

  if (activeUser && users && courses) {
    const course = courses.filter(c => c.id === courseId)[0];
    if (course) {
      const courseGradedUsers = course.enrollments.map(
        e => e.associated_user_id || e.user_id
      );
      if (!courseGradedUsers.includes(activeUserId)) {
        return (
          <div>
            <Typography.Title level={3}>
              {activeUser.name} isn't in {course.name}.
            </Typography.Title>
            <Typography.Text>
              However, the students below are-- click to switch to them:
            </Typography.Text>
            <ul>
              {courseGradedUsers.map(uId => (
                <li key={uId}>
                  <Button
                    type="link"
                    onClick={() => dispatch(changeActiveUser(uId))}
                  >
                    {users[uId].name}
                  </Button>
                </li>
              ))}
            </ul>
          </div>
        );
      }
    }
  }

  if (err) {
    return <ConnectedErrorModal error={err} />;
  }

  if (
    !grades[activeUserId][courseId] ||
    grades[activeUserId][courseId].grade.grade === 'N/A'
  ) {
    return (
      <div align="center">
        <Typography.Title level={3}>
          Grade Breakdown isn't available for this course
          {observees.length > 1 && ' for this student'}.
        </Typography.Title>
        <Link to="/dashboard/grades">
          <Button type="primary">Back to Grades</Button>
        </Link>
      </div>
    );
  }

  if (
    !activeUser ||
    !session ||
    !courses ||
    !allOutcomes ||
    !outcomeResults ||
    !outcomeResults[courseId] ||
    !assignments ||
    !assignments[courseId]
  ) {
    return <Loading text={loadingText} />;
  }

  const course = props.courses.filter(c => c.id === courseId)[0];
  if (!course) {
    notification.error({
      message: 'Unknown Course',
      description: `Couldn't find a course with the specified ID.`
    });
    return <Redirect to="/dashboard/grades" />;
  }

  const averageGrade = gradeAverages ? gradeAverages[courseId] : gradeAverages;
  const grade = grades[activeUserId][courseId];

  const rollupScores = grade.averages;

  function getLowestOutcome() {
    const rollupScore = Object.entries(rollupScores).sort(
      (a, b) => a[1].average - b[1].average
    )[0];
    const outcome = outcomes.filter(o => o.id === parseInt(rollupScore[0]))[0];
    return {
      rollupScore: rollupScore[1],
      outcome
    };
  }

  const lowestOutcome = getLowestOutcome();

  const outcomeTableData = Object.entries(rollupScores)
    .map(([oId, avg]) => {
      const outcome = outcomes.filter(o => o.id === parseInt(oId))[0];
      if (!outcome) {
        return {};
      }

      const results = outcomeResults[courseId][activeUserId][oId];
      const lastAssignmentResult = results.sort((a, b) =>
        dateAsc(a.submitted_or_assessed_at, b.submitted_or_assessed_at)
      )[0];
      const lastAssignmentId = parseInt(
        lastAssignmentResult.links.assignment.split('_')[1]
      );
      const lastAssignment = lastAssignmentResult
        ? assignments[courseId].filter(a => a.id === lastAssignmentId)[0]
        : {};

      // use alignments to figure out things like lastAssignment and timesAssessed
      return {
        name: outcome ? outcome.display_name || outcome.title : 'Error',
        score: +avg.average.toFixed(2),
        worstScoreDropped: avg.did_drop_worst_score,
        lastAssignment: lastAssignment ? lastAssignment.name : 'Unavailable',
        timesAssessed: results.length,
        key: outcome.id,
        id: outcome.id,
        // can be reworked to use the new outcome_alignments
        assignmentTableData: results
          .filter(or => parseInt(or.links.learning_outcome) === outcome.id)
          .map(r => {
            const linkedAssignmentId = parseInt(
              r.links.assignment.split('_')[1]
            );
            const assignment = assignments[courseId].filter(
              a => a.id === linkedAssignmentId
            )[0];
            return {
              assignmentName: assignment ? assignment.name : 'unavailable',
              assignmentUrl: assignment ? assignment.html_url : 'unavailable',
              score: {
                score: r.score,
                possible: r.possible,
                percent: r.percent * 100
              },
              lastSubmission: moment(r.submitted_or_assessed_at).calendar(),
              masteryReached: r.mastery,
              key: linkedAssignmentId
            };
          })
      };
    })
    .filter(otd => !!otd.key);

  const assignmentsByOutcome =
    outcomeAlignments &&
    outcomeAlignments[courseId] &&
    outcomes.reduce((acc = {}, o) => {
      acc[o.id] = flatten(
        outcomeAlignments[courseId]
          .filter(oa => oa.learning_outcome_id === o.id)
          .map(oa =>
            assignments[courseId].filter(a => a.id === oa.assignment_id)
          )
      );

      return acc;
    }, {});

  if (isMobile) {
    return (
      <div>
        <Typography.Title level={2}>Grade Breakdown</Typography.Title>
        <Typography.Text type="secondary">{course.name}</Typography.Text>
        <MobileWhiteSpace />
        <GradeCard
          currentGrade={grade.grade}
          averageGrade={averageGrade}
          userHasValidSubscription={session.has_valid_subscription}
        />
        <MobileWhiteSpace />
        <OutcomeInfo
          lowestOutcome={lowestOutcome}
          min={grade.grade.all_above}
          outcomeRollupScores={rollupScores}
          grade={grade}
          userHasValidSubscription={session.has_valid_subscription}
        />
        <MobileWhiteSpace />
        <Typography.Title level={3}>Outcomes</Typography.Title>
        <MobileAccordion>
          {outcomeTableData.map(d => (
            <MobileAccordion.Panel header={d.name} key={d.key}>
              <MobileList>
                <MobileList.Item multipleLine wrap>
                  <Typography.Text>{d.name}</Typography.Text>
                </MobileList.Item>
                <MobileList.Item extra={d.score}>Score</MobileList.Item>
                <MobileList.Item extra={d.worstScoreDropped ? 'Yes' : 'No'}>
                  Lowest Score Dropped
                </MobileList.Item>
                <MobileList.Item extra={d.timesAssessed}>
                  Times Assessed
                </MobileList.Item>
                <MobileList.Item multipleLine wrap>
                  <Icon component={PlusIcon} /> Average Score <br />
                  <ConnectedAverageOutcomeScore outcomeId={d.id} />
                </MobileList.Item>
                <MobileAccordion>
                  <MobileAccordion.Panel header={<div>Future Assignments</div>}>
                    {
                      <FutureAssignmentsForOutcome
                        outcomeAssignments={assignmentsByOutcome[d.id]}
                      />
                    }
                  </MobileAccordion.Panel>
                </MobileAccordion>
                <MobileAccordion>
                  <MobileAccordion.Panel header="Assignments">
                    {d.assignmentTableData.map(atd => (
                      <MobileAccordion key={atd.key}>
                        <MobileAccordion.Panel
                          header={atd.assignmentName}
                          style={{ paddingLeft: 10 }}
                        >
                          <MobileList style={{ paddingLeft: 10 }}>
                            <MobileList.Item extra={atd.score.score}>
                              Score
                            </MobileList.Item>
                            <MobileList.Item extra={atd.score.possible}>
                              Possible
                            </MobileList.Item>
                            <MobileList.Item extra={atd.lastSubmission}>
                              Last Submission
                            </MobileList.Item>
                            <MobileList.Item>
                              <PopoutLink url={atd.assignmentUrl}>
                                Open on Canvas <Icon component={PopOutIcon} />
                              </PopoutLink>
                            </MobileList.Item>
                          </MobileList>
                        </MobileAccordion.Panel>
                      </MobileAccordion>
                    ))}
                  </MobileAccordion.Panel>
                </MobileAccordion>
              </MobileList>
            </MobileAccordion.Panel>
          ))}
        </MobileAccordion>
      </div>
    );
  }

  return (
    <div>
      <Typography.Title level={2}>
        Grade Breakdown for {course.name}
      </Typography.Title>
      <Row gutter={12}>
        <Col span={8}>
          <GradeCard
            currentGrade={grade.grade}
            averageGrade={averageGrade}
            userHasValidSubscription={session.has_valid_subscription}
          />
        </Col>
        <Col span={16}>
          <OutcomeInfo
            lowestOutcome={lowestOutcome}
            min={grade.grade.all_above}
            outcomeRollupScores={rollupScores}
            grade={grade}
            userHasValidSubscription={session.has_valid_subscription}
          />
        </Col>
      </Row>

      <div style={{ padding: '15px' }} />

      <Typography.Title level={3}>Outcomes</Typography.Title>
      <Table
        columns={outcomeTableColumns}
        dataSource={outcomeTableData}
        expandedRowRender={record =>
          record.assignmentTableData.length > 0 ? (
            <div>
              <Typography.Title level={4}>
                <Icon component={PlusIcon} style={{ paddingRight: '5px' }} />
                Average Score
              </Typography.Title>
              <ConnectedAverageOutcomeScore outcomeId={record.id} />
              <Typography.Title level={4}>Future Assignments</Typography.Title>
              <FutureAssignmentsForOutcome
                outcomeAssignments={assignmentsByOutcome[record.id]}
              />
              <Typography.Title level={4}>Assignments</Typography.Title>
              <Table
                columns={assignmentTableOutcomes}
                dataSource={record.assignmentTableData}
              />
            </div>
          ) : (
            <Typography.Text>
              Couldn't get assignments for this outcome.
            </Typography.Text>
          )
        }
      />
      <Typography.Text type="secondary">
        Please note that these grades may not be accurate or representative of
        your real grade. For the most accurate and up-to-date information,
        please consult someone from your school.
      </Typography.Text>
    </div>
  );
}

const ConnectedGradeBreakdown = connect(state => ({
  loading: state.loading,
  error: state.error,
  courses: state.canvas.courses,
  grades: state.canvas.grades,
  token: state.canvas.token,
  subdomain: state.canvas.subdomain,
  outcomes: state.canvas.outcomes,
  outcomeRollups: state.canvas.outcomeRollups,
  outcomeResults: state.canvas.outcomeResults,
  outcomeAlignments: state.canvas.outcomeAlignments,
  observees: state.canvas.observees,
  assignments: state.canvas.assignments,
  user: state.canvas.user,
  activeUserId: state.canvas.activeUserId,
  users: state.canvas.users,
  session: state.plus.session,
  gradeAverages: state.plus.averages
}))(GradeBreakdown);

export default ConnectedGradeBreakdown;
