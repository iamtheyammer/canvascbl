import React, { useState, useEffect } from 'react';
import { connect } from 'react-redux';
import { Redirect, Link } from 'react-router-dom';
import v4 from 'uuid/v4';
import moment from 'moment';
import { isMobile } from 'react-device-detect';
import { flatten } from 'lodash';

import {
  Typography,
  Spin,
  notification,
  Row,
  Col,
  Button,
  Table,
  Icon
} from 'antd';

import {
  WhiteSpace as MobileWhiteSpace,
  Accordion as MobileAccordion,
  List as MobileList
} from 'antd-mobile';

import './index.css';

import {
  getOutcomeResultsForCourse,
  getUserCourses,
  getAssignmentsForCourse,
  getOutcomeRollupsAndOutcomesForCourse,
  getOutcomeAlignmentsForCourse
} from '../../../../actions/canvas';
import calculateGradeFromOutcomes, {
  gradeMapByGrade
} from '../../../../util/canvas/calculateGradeFromOutcomes';

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
  const [getCoursesId, setGetCoursesId] = useState('');
  const [getRollupsId, setGetRollupsId] = useState('');
  const [getResultsId, setGetResultsId] = useState('');
  const [getAssignmentsId, setGetAssignmentsId] = useState('');
  const [getOutcomeAlignmentsId, setGetOutcomeAlignmentsId] = useState('');
  const [getPlusAverageId, setGetPlusAverageId] = useState('');

  const [loadingText, setLoadingText] = useState('');

  const {
    dispatch,
    token,
    subdomain,
    loading,
    error,
    user,
    courses,
    outcomeRollups,
    outcomeResults,
    outcomeAlignments,
    assignments,
    session,
    gradeAverages
  } = props;

  const err =
    error[getCoursesId] ||
    error[getRollupsId] ||
    error[getResultsId] ||
    error[getAssignmentsId];

  const courseId = parseInt(props.match.params.courseId);

  const allOutcomes = props.outcomes;

  useEffect(
    () => {
      if (isNaN(courseId)) return;

      // loading before fetch because we don't want to request twice
      if (
        loading.includes(getCoursesId) ||
        loading.includes(getRollupsId) ||
        err
      ) {
        return;
      }

      if (!courses && !getCoursesId) {
        const id = v4();
        dispatch(getUserCourses(id, token, subdomain));
        setGetCoursesId(id);
        setLoadingText('your courses');
        return;
      }

      // we can display the page without loading alignments, plus feature
      if (
        !getOutcomeAlignmentsId &&
        session &&
        session.hasValidSubscription &&
        (!outcomeAlignments || !outcomeAlignments[courseId])
      ) {
        const id = v4();
        dispatch(
          getOutcomeAlignmentsForCourse(id, courseId, user.id, token, subdomain)
        );
        setGetOutcomeAlignmentsId(id);
      }

      // we can display the page without the average loading
      if (
        !getPlusAverageId &&
        session &&
        session.hasValidSubscription &&
        (!gradeAverages || !gradeAverages[courseId])
      ) {
        const id = v4();
        dispatch(getAverageGradeForCourse(id, courseId));
        setGetPlusAverageId(id);
      }

      if (
        user &&
        user.id &&
        (!outcomeRollups || !allOutcomes || !allOutcomes[courseId]) &&
        !getRollupsId
      ) {
        const id = v4();
        dispatch(
          getOutcomeRollupsAndOutcomesForCourse(
            id,
            user.id,
            courseId,
            token,
            subdomain
          )
        );
        setLoadingText('your grades');
        setGetRollupsId(id);
      }

      if (
        user &&
        user.id &&
        (!outcomeResults || !outcomeResults[courseId]) &&
        !getResultsId
      ) {
        const id = v4();
        dispatch(
          getOutcomeResultsForCourse(id, user.id, courseId, token, subdomain)
        );
        setLoadingText('your grade in this class');
        setGetResultsId(id);
      }

      if (
        user &&
        user.id &&
        (!assignments || !assignments[courseId]) &&
        !getAssignmentsId
      ) {
        const id = v4();
        dispatch(getAssignmentsForCourse(id, courseId, token, subdomain));
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

  if (err) {
    return <ConnectedErrorModal error={err} />;
  }

  if (
    !user ||
    !session ||
    !courses ||
    !allOutcomes ||
    !allOutcomes[courseId] ||
    !outcomeRollups ||
    !outcomeResults ||
    !outcomeResults[courseId] ||
    !assignments ||
    !assignments[courseId]
  ) {
    return (
      <div align="center">
        <Spin size="default" />
        <span style={{ marginTop: '10px' }} />
        <Typography.Title level={3}>
          {`Loading ${loadingText}...`}
        </Typography.Title>
      </div>
    );
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
  const grade = calculateGradeFromOutcomes({
    [courseId]: props.outcomeRollups[courseId]
  })[courseId];
  const outcomes = props.outcomes[courseId];
  const rollupScores = props.outcomeRollups[courseId][0].scores;

  if (!grade || grade.grade === 'N/A') {
    return (
      <div align="center">
        <Typography.Title level={3}>
          Grade Breakdown Isn't Available for {course.name}.
        </Typography.Title>
        <Link to="/dashboard/grades">
          <Button type="primary">Back to Grades</Button>
        </Link>
      </div>
    );
  }

  const { min } = gradeMapByGrade[grade.grade];

  function getLowestOutcome() {
    const rollupScore = rollupScores.filter(
      rs => rs.score === grade.lowestOutcome
    )[0];
    const outcome = outcomes.filter(
      o => o.id === parseInt(rollupScore.links.outcome)
    )[0];
    return {
      outcome,
      rollupScore
    };
  }

  const lowestOutcome = getLowestOutcome();

  const results = outcomeResults[courseId];

  const outcomeTableData = rollupScores.map(rs => {
    const outcome = outcomes.filter(
      o => o.id === parseInt(rs.links.outcome)
    )[0];
    return {
      name: outcome.display_name || outcome.title,
      score: rs.score,
      lastAssignment: rs.title,
      timesAssessed: rs.count,
      key: outcome.id,
      id: outcome.id,
      assignmentTableData: results
        .filter(or => parseInt(or.links.learning_outcome) === outcome.id)
        .map(r => {
          const linkedAssignmentId = parseInt(r.links.assignment.split('_')[1]);
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
  });

  // only exists if the user has a current session to save the load time
  const assignmentsByOutcome =
    session &&
    session.hasValidSubscription &&
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
          userHasValidSubscription={session.hasValidSubscription}
        />
        <MobileWhiteSpace />
        <OutcomeInfo
          lowestOutcome={lowestOutcome}
          min={min}
          outcomeRollupScores={rollupScores}
          grade={grade}
          userHasValidSubscription={session.hasValidSubscription}
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
                <MobileList.Item extra={d.timesAssessed}>
                  Times Assessed
                </MobileList.Item>
                <MobileList.Item multipleLine wrap>
                  <Icon component={PlusIcon} /> Average Score <br />
                  <ConnectedAverageOutcomeScore outcomeId={d.id} />
                </MobileList.Item>
                <MobileAccordion>
                  <MobileAccordion.Panel
                    header={
                      <div>
                        <PlusIcon style={{ height: '1em' }} /> Future
                        Assignments
                      </div>
                    }
                  >
                    {
                      <FutureAssignmentsForOutcome
                        userHasValidSubscription={session.hasValidSubscription}
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
            userHasValidSubscription={session.hasValidSubscription}
          />
        </Col>
        <Col span={16}>
          <OutcomeInfo
            lowestOutcome={lowestOutcome}
            min={min}
            outcomeRollupScores={rollupScores}
            grade={grade}
            userHasValidSubscription={session.hasValidSubscription}
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
              <Typography.Title level={4}>
                <Icon component={PlusIcon} style={{ paddingRight: '5px' }} />
                Future Assignments
              </Typography.Title>
              <FutureAssignmentsForOutcome
                userHasValidSubscription={session.hasValidSubscription}
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
  token: state.canvas.token,
  subdomain: state.canvas.subdomain,
  outcomes: state.canvas.outcomes,
  outcomeRollups: state.canvas.outcomeRollups,
  outcomeResults: state.canvas.outcomeResults,
  outcomeAlignments: state.canvas.outcomeAlignments,
  assignments: state.canvas.assignments,
  user: state.canvas.user,
  session: state.plus.session,
  gradeAverages: state.plus.averages
}))(GradeBreakdown);

export default ConnectedGradeBreakdown;
