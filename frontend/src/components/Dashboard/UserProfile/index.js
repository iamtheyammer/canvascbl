import React, { Fragment } from 'react';
import { connect } from 'react-redux';
import { isMobile } from 'react-device-detect';
import v4 from 'uuid/v4';
import moment from 'moment';

import {
  Descriptions,
  Typography,
  List,
  Skeleton,
  Button,
  Popconfirm
} from 'antd';
import {
  Card as MobileCard,
  List as MobileList,
  Accordion as MobileAccordion,
  Button as MobileButton
} from 'antd-mobile';
import Padding from '../../Padding';
import { getAuthorizedApps, revokeGrant } from '../../../actions/oauth2';
import { setGetAuthorizedAppsId } from '../../../actions/components/userprofile';

function UserProfile(props) {
  const {
    user,
    observees,
    loading,
    getAuthorizedAppsId,
    authorizedApps,
    getAuthorizedAppsError,
    revokedGrantIds,
    revokeGrantErrors,
    dispatch
  } = props;

  const loadingApps = !authorizedApps || loading.includes(getAuthorizedAppsId);
  if (!getAuthorizedAppsId && !authorizedApps) {
    const id = v4();
    dispatch(getAuthorizedApps(id));
    dispatch(setGetAuthorizedAppsId(id));
  }

  const credsMap =
    authorizedApps &&
    authorizedApps.credentials.reduce((acc, val) => {
      acc[val.id] = val.name;
      return acc;
    }, {});

  if (isMobile) {
    return (
      <>
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
        {observees && observees.length > 0 && (
          <>
            <Padding br />
            <Typography.Title level={3}>My Observees</Typography.Title>
            <MobileList>
              {observees.map(o => (
                <MobileList.Item key={o.id}>{o.name}</MobileList.Item>
              ))}
            </MobileList>
          </>
        )}
        <Typography.Title level={3}>Authorized Apps</Typography.Title>
        <Typography.Text type="secondary">
          You granted the apps below access to your CanvasCBL account.
        </Typography.Text>
        <Padding all={5} />
        {authorizedApps && !authorizedApps.grants ? (
          'No apps granted access.'
        ) : !getAuthorizedAppsError ? (
          !loadingApps ? (
            <MobileAccordion>
              {authorizedApps.grants.map(grant =>
                revokedGrantIds && revokedGrantIds.includes(grant.id) ? (
                  <Fragment key={grant.id} />
                ) : (
                  <MobileAccordion.Panel
                    header={credsMap[grant.oauth2_credential_id]}
                    key={grant.id}
                  >
                    <MobileList>
                      <MobileList.Item>
                        Authorized at/on:
                        <MobileList.Item.Brief>
                          {`${moment(grant.inserted_at).calendar()}${
                            revokeGrantErrors && revokeGrantErrors[grant.id]
                              ? ' | Error revoking. Try again later.'
                              : ''
                          }`}
                        </MobileList.Item.Brief>
                      </MobileList.Item>
                      <MobileList.Item>
                        <MobileButton
                          disabled={
                            revokeGrantErrors && revokeGrantErrors[grant.id]
                          }
                          onClick={() =>
                            dispatch(
                              revokeGrant(/*id doesn't matter*/ v4(), grant.id)
                            )
                          }
                        >
                          Revoke
                        </MobileButton>
                      </MobileList.Item>
                    </MobileList>
                  </MobileAccordion.Panel>
                )
              )}
            </MobileAccordion>
          ) : (
            <Skeleton active />
          )
        ) : (
          <Typography.Text type="danger">
            Error getting authorized apps.
          </Typography.Text>
        )}
      </>
    );
  }

  return (
    <>
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
        <>
          <Padding br />
          <Typography.Title level={3}>My Observees</Typography.Title>
          <Typography.Text type="secondary">
            Switch between observees using the menu at the top right.
          </Typography.Text>
          <Padding br />
          <ul>
            {observees.map(o => (
              <li key={o.id}>{o.name}</li>
            ))}
          </ul>
        </>
      )}
      <Typography.Title level={3}>Authorized Apps</Typography.Title>
      <Typography.Text type="secondary">
        You granted the apps below access to your CanvasCBL account.
      </Typography.Text>
      <Padding all={10} />
      {authorizedApps && !authorizedApps.grants ? (
        'No apps granted access.'
      ) : !getAuthorizedAppsError ? (
        !loadingApps ? (
          <List
            itemLayout="horizontal"
            dataSource={authorizedApps.grants}
            renderItem={grant =>
              revokedGrantIds && revokedGrantIds.includes(grant.id) ? (
                <Fragment key={grant.id} />
              ) : (
                <List.Item
                  actions={[
                    <Popconfirm
                      key="revoke"
                      placement="topRight"
                      title={
                        'Are you sure you want to revoke this authorization?'
                      }
                      onConfirm={() =>
                        dispatch(
                          revokeGrant(/*id doesn't matter*/ v4(), grant.id)
                        )
                      }
                    >
                      <Button
                        type="link"
                        disabled={
                          revokeGrantErrors && revokeGrantErrors[grant.id]
                        }
                      >
                        Revoke
                      </Button>
                    </Popconfirm>
                  ]}
                >
                  <List.Item.Meta
                    title={credsMap[grant.oauth2_credential_id]}
                    description={`Authorized at/on: ${moment(
                      grant.inserted_at
                    ).calendar()}${
                      revokeGrantErrors && revokeGrantErrors[grant.id]
                        ? ' | Error revoking. Try again later.'
                        : ''
                    }`}
                  />
                </List.Item>
              )
            }
          />
        ) : (
          <Skeleton active />
        )
      ) : (
        <Typography.Text type="danger">
          Error getting authorized apps.
        </Typography.Text>
      )}
    </>
  );
}

const ConnectedUserProfile = connect(state => ({
  user: state.canvas.user,
  observees: state.canvas.observees,
  loading: state.loading,
  getAuthorizedAppsId: state.components.userprofile.getAuthorizedAppsId,
  authorizedApps: state.oauth2.authorizedApps,
  getAuthorizedAppsError: state.oauth2.getAuthorizedAppsError,
  revokedGrantIds: state.oauth2.revokedGrantIds,
  revokeGrantErrors: state.oauth2.revokeGrantErrors
}))(UserProfile);

export default ConnectedUserProfile;
