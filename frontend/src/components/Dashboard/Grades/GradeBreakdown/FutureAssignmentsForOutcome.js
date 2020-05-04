import React from 'react';
import * as PropTypes from 'prop-types';
import moment from 'moment';
import { Divider, Icon, Table, Typography } from 'antd';
import { Accordion as MobileAccordion, List as MobileList } from 'antd-mobile';
import { isMobile } from 'react-device-detect';
import PopoutLink from '../../../PopoutLink';
import { dateDesc } from '../../../../util/sort';
import v4 from 'uuid/v4';
import * as sort from '../../../../util/sort';
import { ReactComponent as PopOutIcon } from '../../../../assets/pop_out.svg';
import {
  destinationNames,
  destinationTypes,
  vias
} from '../../../../util/tracking';

const tableColumns = [
  {
    title: 'Assignment Name',
    dataIndex: 'name'
  },
  {
    title: 'Due Date',
    dataIndex: 'dueDate',
    sorter: (a, b) => dateDesc(a.dueAt, b.dueAt)
  },
  {
    title: 'Actions',
    dataIndex: 'actions',
    render: (text, record) =>
      record.actions.map((a, i) => (
        <div key={i}>
          {a}
          {i + 1 !== record.actions.length && <Divider type="vertical" />}
        </div>
      )),
    mobileRender: (text, record) => (
      <MobileList>
        {record.actions.map((a) => (
          <MobileList.Item key={v4()}>{a}</MobileList.Item>
        ))}
      </MobileList>
    )
  }
];

function FutureAssignmentsForOutcome(props) {
  const { outcomeAssignments } = props;

  const futureAssignments = outcomeAssignments.filter(
    (a) => a.due_at && moment(a.due_at).isAfter()
  );

  const data = futureAssignments.map((fa) => ({
    name: fa.name,
    dueDate: moment(fa.due_at).calendar(),
    dueAt: fa.due_at,
    url: fa.html_url,
    actions: [
      <PopoutLink
        url={fa.html_url}
        addIcon
        tracking={{
          destinationName: destinationNames.canvas,
          destinationType: destinationTypes.assignment,
          via:
            vias.gradeBreakdownOutcomesTableFutureAssignmentsTableOpenOnCanvas
        }}
      >
        Open on Canvas
      </PopoutLink>
    ],
    key: fa.id
  }));

  if (data.length === 0) {
    return (
      <Typography.Text>
        There are no future posted assignments for this outcome.
      </Typography.Text>
    );
  }

  if (isMobile) {
    return data.sort(sort.dateAsc).map((a) => (
      <MobileAccordion key={a.key}>
        <MobileAccordion.Panel header={a.name} style={{ paddingLeft: 10 }}>
          <MobileList style={{ paddingLeft: 10 }}>
            <MobileList.Item extra={a.dueDate}>Due Date</MobileList.Item>
            <MobileList.Item>
              <PopoutLink
                url={a.url}
                tracking={{
                  destinationName: destinationNames.canvas,
                  destinationType: destinationTypes.assignment,
                  via:
                    vias.gradeBreakdownOutcomesTableFutureAssignmentsTableOpenOnCanvas
                }}
              >
                Open on Canvas <Icon component={PopOutIcon} />
              </PopoutLink>
            </MobileList.Item>
          </MobileList>
        </MobileAccordion.Panel>
      </MobileAccordion>
    ));
  }

  return <Table columns={tableColumns} dataSource={data} />;
}

FutureAssignmentsForOutcome.propTypes = {
  // assignments that map to this outcome
  outcomeAssignments: PropTypes.arrayOf(
    PropTypes.shape({
      due_at: PropTypes.string,
      html_url: PropTypes.string,
      name: PropTypes.string
    })
  )
};

export default FutureAssignmentsForOutcome;
