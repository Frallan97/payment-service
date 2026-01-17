import { useEffect, useState } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';
import { AuthService } from '../lib/auth';

export default function OAuthCallback() {
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const accessToken = searchParams.get('access_token');
    const errorParam = searchParams.get('error');

    if (errorParam) {
      setError(errorParam);
      setTimeout(() => navigate('/login'), 3000);
      return;
    }

    if (accessToken) {
      // Store the token
      AuthService.setToken(accessToken);
      // Redirect to dashboard
      navigate('/dashboard');
    } else {
      setError('No access token received');
      setTimeout(() => navigate('/login'), 3000);
    }
  }, [searchParams, navigate]);

  return (
    <div className="min-h-screen flex items-center justify-center bg-background">
      <div className="text-center">
        {error ? (
          <>
            <div className="text-red-500 text-xl mb-4">Authentication Error</div>
            <p className="text-muted-foreground">{error}</p>
            <p className="text-sm text-muted-foreground mt-2">Redirecting to login...</p>
          </>
        ) : (
          <>
            <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary mx-auto mb-4"></div>
            <p className="text-muted-foreground">Completing authentication...</p>
          </>
        )}
      </div>
    </div>
  );
}
