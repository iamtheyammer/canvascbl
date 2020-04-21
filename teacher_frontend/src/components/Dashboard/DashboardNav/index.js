import React from "react";
import { Layout, Menu } from "antd";
import styled from "styled-components";
import logo from "../../../assets/banner-light.svg";

const { Header } = Layout;

const StyledLogo = styled.img`
  width: 126px;
  height: auto;
  margin: 13px 18px 0 0;
  float: left;
`;

const StyledMenu = styled(Menu)`
  line-height: 64px;
  float: left;
`;

const StyledRightMenu = styled(Menu)`
  line-height: 64px;
  float: right;
`;

function DashboardNav(props) {
  return (
    <Header>
      <StyledLogo src={logo} alt="CanvasCBL Logo" />
      <StyledMenu
        theme="dark"
        mode="horizontal"
        defaultSelectedKeys="/dashboard/courses"
      >
        <Menu.Item key="/dashboard/courses">Courses</Menu.Item>
      </StyledMenu>
      <StyledRightMenu theme="dark" mode="horizontal">
        <Menu.Item key="/dashboard/logout">Logout</Menu.Item>
      </StyledRightMenu>
    </Header>
  );
}

export default DashboardNav;
