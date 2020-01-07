import React from 'react';
import { connect } from 'react-redux';
import { isMobile } from 'react-device-detect';

import { Spin, Descriptions, Typography } from 'antd';
import { Card as MobileCard, List as MobileList } from 'antd-mobile';
import Padding from '../../Padding';
import { Link } from 'react-router-dom';

function UserProfile(props) {
  const { user, observees } = props;

  if (!user) {
    return (
      <div align="center" style={{ marginTop: '20px' }}>
        <Spin size="large" />
      </div>
    );
  }

  if (isMobile) {
    return (
      <div>
        <Typography.Title level={2}>My Profile</Typography.Title>
        <MobileCard>
          <MobileCard.Header title={user.name} thumb={user.avatar_url} />
          <MobileCard.Body>
            Email: {user.primary_email}
            <br />
            Time Zone: {user.time_zone}
          </MobileCard.Body>
          <MobileCard.Footer content={`ID: ${user.id}`} />
        </MobileCard>
        {observees && observees.length > 1 && (
          <div>
            <Padding br />
            <Typography.Title level={3}>My Observees</Typography.Title>
            <MobileList>
              {observees.map(o => (
                <MobileList.Item key={o.id}>{o.name}</MobileList.Item>
              ))}
            </MobileList>
          </div>
        )}
      </div>
    );
  }

  return (
    <div>
      <Typography.Title level={2}>My Profile</Typography.Title>
      <Descriptions bordered>
        <Descriptions.Item label="Name">{user.name}</Descriptions.Item>
        <Descriptions.Item label="Email">
          {user.primary_email}
        </Descriptions.Item>
        <Descriptions.Item label="ID">{user.id}</Descriptions.Item>
        <Descriptions.Item label="Time Zone">
          {user.time_zone}
        </Descriptions.Item>
      </Descriptions>
      {observees && observees.length > 0 && (
        <div>
          <Padding br />
          <Typography.Title level={3}>My Observees</Typography.Title>
          <Typography.Text type="secondary">
            Switch between observees using the menu at the top right. (Not
            there? Click <Link to="/dashboard/grades">here</Link>)
          </Typography.Text>
          <Padding br />
          <ul>
            {observees.map(o => (
              <li key={o.id}>{o.name}</li>
            ))}
          </ul>
        </div>
      )}
    </div>
  );
}

const ConnectedUserProfile = connect(state => ({
  user: state.canvas.user,
  observees: state.canvas.observees
}))(UserProfile);

export default ConnectedUserProfile;
