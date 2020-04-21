import React, { Fragment } from 'react';
import { Col, Row, Skeleton, Typography } from 'antd';
import { connect } from 'react-redux';
import styled from 'styled-components';
import { Card } from 'antd';
import Padding from '../../Padding';
import { chunk } from 'lodash';

const StyledCard = styled(Card)`
  width: 240px;
`;

function Courses(props) {
  const { courses, distanceLearningPairs, history } = props;

  return (
    <>
      <Typography.Title level={2}>Courses</Typography.Title>
      <Typography.Text type="secondary">
        Click a course to see more.
      </Typography.Text>
      <Padding all={10} />
      {!distanceLearningPairs ? (
        <StyledCard hoverable cover={<Skeleton active />}>
          <Skeleton active />
        </StyledCard>
      ) : (
        chunk(distanceLearningPairs, 4).map((ch, i) => (
          <Fragment key={i}>
            <Row type="flex" justify="start" gutter={24}>
              {ch.map((p) => {
                const dlCourse = courses.filter(
                  (c) => c.id === p.distance_learning_course_id
                )[0];
                const oriCourse = courses.filter(
                  (c) => c.id === p.original_course_id
                )[0];

                const coverImgSrc =
                  dlCourse.image_download_url || oriCourse.image_download_url;

                if (!dlCourse || !oriCourse) {
                  return (
                    <Card
                      key={`${dlCourse.id}_${oriCourse.id}`}
                      hoverable
                      // cover={
                      //   <img
                      //     alt="example"
                      //     src="https://os.alipayobjects.com/rmsportal/QBnOOoLaAfKPirc.png"
                      //   />
                      // }
                    >
                      <Card.Meta
                        title={`Error loading ${p.course_name}`}
                        description={
                          "Weirdly enough, we're missing a course or two."
                        }
                      />
                    </Card>
                  );
                }

                return (
                  <Col span={6} key={`${dlCourse.id}_${oriCourse.id}`}>
                    <Card
                      hoverable
                      onClick={() =>
                        history.push(
                          `/dashboard/courses/${oriCourse.id}_${dlCourse.id}/overview`
                        )
                      }
                      cover={
                        coverImgSrc && (
                          <img alt={`${p.course_name}`} src={coverImgSrc} />
                        )
                      }
                    >
                      <Card.Meta
                        title={p.course_name}
                        description={
                          <Typography.Text>
                            Compiled from {dlCourse.name} and {oriCourse.name}.
                          </Typography.Text>
                        }
                      />
                    </Card>
                  </Col>
                );
              })}
            </Row>
            <Padding all={12} />
          </Fragment>
        ))
      )}
    </>
  );
}

const ConnectedCourses = connect((state) => ({
  coursesAreLoading: state.canvas.coursesAreLoading,
  courses: state.canvas.courses,
  distanceLearningPairs: state.canvas.distanceLearningPairs,
}))(Courses);

export default ConnectedCourses;
