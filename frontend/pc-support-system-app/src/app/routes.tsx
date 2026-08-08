import { BrowserRouter, Routes, Route } from "react-router-dom";
import Home from "../features/home/pages/Home";
import Register from "../features/auth/pages/Register";
import Login from "../features/auth/pages/Login";
import ProtectedRoute from "../routes/ProtectedRoute";
import CustomerDashboard from "../features/customer/pages/Dashboard";
import Profile from "../features/customer/pages/Profile";
import Devices from "../features/customer/pages/Devices";
import JobDetails from "../features/customer/pages/JobDetails";
import StaffLayout from "../layouts/StaffLayout";
import Jobs from "../features/customer/pages/Jobs";
import Dashboard from "../features/dashboard/Dashboard";
import StaffProfile from "../features/profile/StaffProfile";
import CustomerManagement from "../features/customer/CustomerManagement";
import CustomerDetails    from "../features/customer/CustomerDetails";
import CustomerForm    from "../features/customer/CustomerForm";
import StaffManagement from "../features/staff/StaffManagement";
import StaffDetails from "../features/staff/StaffDetails";
import DeviceForm from "../features/devices/DeviceForm";
import DeviceDetails from "../features/devices/DeviceDetails";
import AuditLogs from "../features/auditlogs/AuditLogs";
import JobManagement from "../features/jobs/JobManagement";

export default function AppRoutes() {
  return (
    <BrowserRouter>

      <Routes>

        <Route
          path="/"
          element={<Home />}
        />

        <Route
          path="/login"
          element={<Login />}
        />

        <Route
          path="/register"
          element={<Register />}
        />

        <Route
          path="/customer"
          element={
            <ProtectedRoute>
              <CustomerDashboard />
            </ProtectedRoute>
          }
        />

        <Route
          path="/customer/profile"
          element={
            <ProtectedRoute>
              <Profile />
            </ProtectedRoute>
          }
        />
        <Route
    path="/customer/devices"
    element={
        <ProtectedRoute>
            <Devices />
        </ProtectedRoute>
    }
        />
      <Route
    path="/customer/jobs"
    element={
        <ProtectedRoute>
            <Jobs />
        </ProtectedRoute>
    }
/>
      <Route
          path="/customer/jobs/:jobNumber"
          element={
              <ProtectedRoute>
                  <JobDetails />
              </ProtectedRoute>
          }
/>
  <Route
    path="/receptionist"
    element={
        <ProtectedRoute>
            <StaffLayout />
        </ProtectedRoute>
    }
>

    <Route
        index
        element={<Dashboard />}
    />

    <Route
        path="customers"
        element={<CustomerManagement />}
    />

    <Route
    path="customers/create"
    element={<CustomerForm />}
/>

      <Route
    path="customers/:customerId/edit"
    element={<CustomerForm />}
/>

    <Route
    path="customers/:customerId"
    element={<CustomerDetails />}
/>

    <Route
    path="jobs"
    element={<JobManagement />}
/>

    <Route
        path="profile"
        element={<StaffProfile />}
    />

            <Route
            path="customers/:customerId/devices/create"
            element={<DeviceForm />}
        />

        <Route
            path="customers/:customerId/devices/:deviceId"
            element={<DeviceDetails />}
        />

        <Route
            path="customers/:customerId/devices/:deviceId/edit"
            element={<DeviceForm />}
        />

    </Route>

<Route
    path="/technician"
    element={
        <ProtectedRoute>
            <StaffLayout />
        </ProtectedRoute>
    }
>
    <Route
        index
        element={<Dashboard />}
    />

    <Route
        path="customers"
        element={<CustomerManagement />}
    />

    <Route
        path="customers/:customerId"
        element={<CustomerDetails />}
    />

    <Route
    path="jobs"
    element={<JobManagement />}
/>

    <Route
        path="profile"
        element={<StaffProfile />}
    />

    <Route
    path="customers/:customerId/devices/:deviceId"
    element={<DeviceDetails />}
/>

<Route
    path="customers/:customerId/devices/:deviceId/edit"
    element={<DeviceForm />}
/>

</Route>

<Route
    path="/head-technician"
    element={
        <ProtectedRoute>
            <StaffLayout />
        </ProtectedRoute>
    }
>
    <Route
        index
        element={<Dashboard />}
    />

    <Route
        path="customers"
        element={<CustomerManagement />}
    />

    <Route
        path="customers/create"
        element={<CustomerForm />}
    />

    <Route
        path="customers/:customerId/edit"
        element={<CustomerForm />}
    />

    <Route
        path="customers/:customerId"
        element={<CustomerDetails />}
    />

    <Route
        path="staff"
        element={<StaffManagement />}
    />

    <Route
        path="staff/:staffId"
        element={<StaffDetails />}
    />

   <Route
    path="jobs"
    element={<JobManagement />}
/>

    <Route
        path="profile"
        element={<StaffProfile />}
    />

        <Route
            path="customers/:customerId/devices/create"
            element={<DeviceForm />}
        />

        <Route
            path="customers/:customerId/devices/:deviceId"
            element={<DeviceDetails />}
        />

        <Route
            path="customers/:customerId/devices/:deviceId/edit"
            element={<DeviceForm />}
        />

</Route>

<Route
    path="/admin"
    element={
        <ProtectedRoute>
            <StaffLayout />
        </ProtectedRoute>
    }
>
    <Route
        index
        element={<Dashboard />}
    />

    <Route
        path="customers"
        element={<CustomerManagement />}
    />

    <Route
        path="customers/create"
        element={<CustomerForm />}
    />

    <Route
        path="customers/:customerId/edit"
        element={<CustomerForm />}
    />

    <Route
        path="customers/:customerId"
        element={<CustomerDetails />}
    />

    <Route
        path="staff"
        element={<StaffManagement />}
    />

    <Route
        path="staff/:staffId"
        element={<StaffDetails />}
    />

   <Route
    path="jobs"
    element={<JobManagement />}
/>

    <Route
        path="profile"
        element={<StaffProfile />}
    />

    <Route
        path="customers/:customerId/devices/create"
        element={<DeviceForm />}
    />

    <Route
        path="customers/:customerId/devices/:deviceId"
        element={<DeviceDetails />}
    />

    <Route
        path="customers/:customerId/devices/:deviceId/edit"
        element={<DeviceForm />}
    />
    <Route
    path="audit-logs"
    element={<AuditLogs />}
/>

</Route>

 <Route
    path="/super-admin"
    element={
        <ProtectedRoute>
            <StaffLayout />
        </ProtectedRoute>
    }
>
    <Route
        index
        element={<Dashboard />}
    />

    <Route
        path="customers"
        element={<CustomerManagement />}
    />

    <Route
        path="customers/create"
        element={<CustomerForm />}
    />

    <Route
        path="customers/:customerId/edit"
        element={<CustomerForm />}
    />

    <Route
        path="customers/:customerId"
        element={<CustomerDetails />}
    />

    <Route
        path="staff"
        element={<StaffManagement />}
    />

    <Route
        path="staff/:staffId"
        element={<StaffDetails />}
    />

    <Route
    path="jobs"
    element={<JobManagement />}
/>

    <Route
        path="profile"
        element={<StaffProfile />}
    />

        <Route
            path="customers/:customerId/devices/create"
            element={<DeviceForm />}
        />

        <Route
            path="customers/:customerId/devices/:deviceId"
            element={<DeviceDetails />}
        />

        <Route
            path="customers/:customerId/devices/:deviceId/edit"
            element={<DeviceForm />}
        />

        <Route
    path="audit-logs"
    element={<AuditLogs />}
/>
    </Route>

        </Routes>

    </BrowserRouter>
  );
}