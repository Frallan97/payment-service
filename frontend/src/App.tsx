import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom";
import { AuthService } from "./lib/auth";
import DashboardLayout from "./components/admin/DashboardLayout";
import LoginPage from "./pages/LoginPage";
import OAuthCallback from "./components/auth/OAuthCallback";
import DashboardPage from "./pages/DashboardPage";
import PaymentsPage from "./pages/PaymentsPage";
import PaymentDetailPage from "./pages/PaymentDetailPage";
import SubscriptionsPage from "./pages/SubscriptionsPage";
import SubscriptionDetailPage from "./pages/SubscriptionDetailPage";
import RefundsPage from "./pages/RefundsPage";
import { ThemeProvider } from "./components/theme-provider";

function ProtectedRoute({ children }: { children: React.ReactNode }) {
  return AuthService.isAuthenticated() ? (
    <>{children}</>
  ) : (
    <Navigate to="/login" replace />
  );
}

function App() {
  return (
    <ThemeProvider attribute="class" defaultTheme="system" enableSystem>
      <BrowserRouter>
        <Routes>
        <Route path="/login" element={<LoginPage />} />
        <Route path="/auth/callback" element={<OAuthCallback />} />
        <Route
          path="/dashboard"
          element={
            <ProtectedRoute>
              <DashboardLayout />
            </ProtectedRoute>
          }
        >
          <Route index element={<DashboardPage />} />
          <Route path="payments" element={<PaymentsPage />} />
          <Route path="payments/:id" element={<PaymentDetailPage />} />
          <Route path="subscriptions" element={<SubscriptionsPage />} />
          <Route path="subscriptions/:id" element={<SubscriptionDetailPage />} />
          <Route path="refunds" element={<RefundsPage />} />
        </Route>
          <Route path="/" element={<Navigate to="/dashboard" replace />} />
        </Routes>
      </BrowserRouter>
    </ThemeProvider>
  );
}

export default App;
