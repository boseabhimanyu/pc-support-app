import { BrowserRouter, Routes, Route } from "react-router-dom";
import Home from "../features/home/pages/Home";
import Register from "../features/auth/pages/Register";
import Login from "../features/auth/pages/Login";
import ProtectedRoute from "../routes/ProtectedRoute";
import CustomerDashboard from "../features/customer/pages/Dashboard";
import Profile from "../features/customer/pages/Profile";
import Devices from "../features/customer/pages/Devices";
import JobDetails from "../features/customer/pages/JobDetails";

import Jobs from "../features/customer/pages/Jobs";

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

      </Routes>

    </BrowserRouter>
  );
}