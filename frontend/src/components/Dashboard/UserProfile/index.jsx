import React from 'react';
import { connect } from 'react-redux';
import { isMobile } from 'react-device-detect';

import { Spin, Descriptions, Typography } from 'antd';
import { Card as MobileCard } from 'antd-mobile';

function UserProfile(props) {
  const { user } = props;

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
      </div>
    );
  }

  return (
    <Descriptions title="My Profile" bordered>
      <Descriptions.Item label="Name">{user.name}</Descriptions.Item>
      <Descriptions.Item label="Email">{user.primary_email}</Descriptions.Item>
      <Descriptions.Item label="ID">{user.id}</Descriptions.Item>
      <Descriptions.Item label="Time Zone">{user.time_zone}</Descriptions.Item>
    </Descriptions>
  );
}

const ConnectedUserProfile = connect(state => ({
  user: state.canvas.user
}))(UserProfile);

export default ConnectedUserProfile;
