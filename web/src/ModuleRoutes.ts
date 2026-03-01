/*
 * Copyright 2023 Harness, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

export interface ModuleProps {
  space?: string
}

export const modulePathProps = {
  space: ':space*'
}

export interface ModuleRoutes {
  toFeatureFlags: (args: Required<Pick<ModuleProps, 'space'>>) => string
  toTechDebt: (args: Required<Pick<ModuleProps, 'space'>>) => string
  toSecurityScans: (args: Required<Pick<ModuleProps, 'space'>>) => string
  toMonitors: (args: Required<Pick<ModuleProps, 'space'>>) => string
  toErrors: (args: Required<Pick<ModuleProps, 'space'>>) => string
  toQualityGates: (args: Required<Pick<ModuleProps, 'space'>>) => string
}

export const moduleRoutes: ModuleRoutes = {
  toFeatureFlags: ({ space }) => `/spaces/${space}/featureflags`,
  toTechDebt: ({ space }) => `/spaces/${space}/techdebt`,
  toSecurityScans: ({ space }) => `/spaces/${space}/security`,
  toMonitors: ({ space }) => `/spaces/${space}/monitors`,
  toErrors: ({ space }) => `/spaces/${space}/errors`,
  toQualityGates: ({ space }) => `/spaces/${space}/quality`
}
