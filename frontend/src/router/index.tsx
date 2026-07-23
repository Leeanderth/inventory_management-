import { createBrowserRouter } from 'react-router-dom'
import App from '@/App'
import AppLayout from '@/layouts/AppLayout'
import LoginView from '@/views/LoginView'
import ProductView from '@/views/ProductView'
import RolePermissionView from '@/views/RolePermissionView'
import UserManagementView from '@/views/UserManagementView'
import RequirePermissionRoute from './RequirePermissionRoute'
import HomeRedirect from './HomeRedirect'

export const router = createBrowserRouter([
  {
    path: '/',
    element: <App />,
    children: [
      { index: true, element: <HomeRedirect /> },
      { path: 'login', element: <LoginView /> },
      {
        element: <AppLayout />,
        children: [
          {
            path: 'stock',
            element: (
              <RequirePermissionRoute permission="stock:view">
                <ProductView />
              </RequirePermissionRoute>
            ),
          },
          {
            path: 'users',
            element: (
              <RequirePermissionRoute permission="user:view">
                <UserManagementView />
              </RequirePermissionRoute>
            ),
          },
          {
            path: 'roles',
            element: (
              <RequirePermissionRoute permission="role:view">
                <RolePermissionView />
              </RequirePermissionRoute>
            ),
          },
        ],
      },
      { path: '*', element: <HomeRedirect /> },
    ],
  },
])
