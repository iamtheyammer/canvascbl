const {
  override,
  fixBabelImports,
  addWebpackPlugin,
  addBabelPlugin
} = require('customize-cra');
const Dotenv = require('dotenv-webpack');

const isDev = process.env.NODE_ENV === 'development';

module.exports = override(
  fixBabelImports('import-antd', {
    libraryName: 'antd',
    libraryDirectory: 'es',
    style: 'css'
  }),
  fixBabelImports('import-mobile', {
    libraryName: 'antd-mobile',
    style: 'css'
  }),
  addBabelPlugin([
    'babel-plugin-styled-components',
    { ssr: false, displayName: isDev }
  ]),
  // fixBabelImports("import-lodash", {
  //   libraryName: "lodash",
  //   libraryDirectory: "",
  //   camel2DashComponentName: false // default: true
  // }),
  addWebpackPlugin(new Dotenv())
);
