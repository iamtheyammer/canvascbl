import styled from 'styled-components';
import { isMobile } from 'react-device-detect';
import { Card } from 'antd';
import { LoadingOutlined } from '@ant-design/icons';
import banner from '../../assets/banner.svg';
import React from 'react';
import Padding from '../Padding';

const StyledCard = styled(Card)`
  width: 35%;
  margin: 0;
  position: fixed;
  top: 40%;
  left: 50%;
  transform: translateY(-50%) translateX(-50%);
  box-shadow: 2px 2px 15px 0 rgba(255, 255, 255, 0.3);
  .ant-card-head-title {
    padding: 25px !important;
  }
`;

const StyledLogo = styled.div`
  background-image: url("${(props) => props.src}");
  background-position: center;
  background-repeat: no-repeat;
  background-size: contain;
  width: auto;
  height: 80px;
  margin: 0 auto;
  transform: scale(1.1);
`;

const CenterText = styled.div`
  text-align: center;
`;

const LoadingIcon = styled(LoadingOutlined)`
  display: flex;
  justify-content: center;
  align-items: center;
  font-size: 3em;
`;

function MainCard(props) {
  const { children, loading, loadingText } = props;

  if (isMobile) {
    return (
      <Card title={<StyledLogo src={banner} />}>
        {loading ? (
          <>
            <LoadingIcon />
            {loadingText && (
              <>
                <Padding all={5} />
                <CenterText>{loadingText}</CenterText>
              </>
            )}
          </>
        ) : (
          children
        )}
      </Card>
    );
  }

  return (
    <StyledCard title={<StyledLogo src={banner} />}>
      {loading ? (
        <>
          <LoadingIcon />
          {loadingText && (
            <>
              <Padding all={5} />
              <CenterText>{loadingText}</CenterText>
            </>
          )}
        </>
      ) : (
        children
      )}
    </StyledCard>
  );
}

export default MainCard;
