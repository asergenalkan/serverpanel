import { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import { useAuth } from '@/contexts/AuthContext';
import { dashboardAPI } from '@/lib/api';
import { Card, CardContent } from '@/components/ui/Card';
import { Button } from '@/components/ui/Button';
import Layout from '@/components/Layout';
import {
  Users,
  Globe,
  Database,
  Mail,
  Cpu,
  HardDrive,
  MemoryStick,
} from 'lucide-react';

interface DashboardStats {
  total_users: number;
  total_domains: number;
  total_databases: number;
  total_emails: number;
  system_stats?: {
    cpu_usage: number;
    memory_total: number;
    memory_used: number;
    disk_total: number;
    disk_used: number;
    load_average: number[];
  };
}

function formatBytes(bytes: number): string {
  if (bytes === 0) return '0 B';
  const k = 1024;
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i];
}

function StatCard({
  title,
  value,
  icon: Icon,
  color,
}: {
  title: string;
  value: string | number;
  icon: React.ComponentType<{ className?: string }>;
  color: string;
}) {
  return (
    <Card className="hover:shadow-lg transition-shadow">
      <CardContent className="p-6">
        <div className="flex items-center justify-between">
          <div>
            <p className="text-sm font-medium text-muted-foreground">{title}</p>
            <p className="text-2xl font-bold mt-1">{value}</p>
          </div>
          <div className={`p-3 rounded-xl ${color}`}>
            <Icon className="w-6 h-6 text-white" />
          </div>
        </div>
      </CardContent>
    </Card>
  );
}

function SystemStatCard({
  title,
  used,
  total,
  icon: Icon,
}: {
  title: string;
  used: number;
  total: number;
  icon: React.ComponentType<{ className?: string }>;
}) {
  const percentage = total > 0 ? Math.round((used / total) * 100) : 0;
  const color = percentage > 80 ? 'bg-red-500' : percentage > 60 ? 'bg-yellow-500' : 'bg-green-500';

  return (
    <Card>
      <CardContent className="p-6">
        <div className="flex items-center gap-4">
          <div className="p-3 rounded-xl bg-slate-100">
            <Icon className="w-6 h-6 text-slate-600" />
          </div>
          <div className="flex-1">
            <div className="flex justify-between items-center mb-2">
              <p className="text-sm font-medium">{title}</p>
              <span className="text-sm text-muted-foreground">{percentage}%</span>
            </div>
            <div className="h-2 bg-slate-100 rounded-full overflow-hidden">
              <div className={`h-full ${color} transition-all`} style={{ width: `${percentage}%` }} />
            </div>
            <p className="text-xs text-muted-foreground mt-1">
              {formatBytes(used)} / {formatBytes(total)}
            </p>
          </div>
        </div>
      </CardContent>
    </Card>
  );
}

export default function Dashboard() {
  const { user } = useAuth();
  const [stats, setStats] = useState<DashboardStats | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchStats = async () => {
      try {
        const response = await dashboardAPI.getStats();
        if (response.data.success) {
          setStats(response.data.data);
        }
      } catch (error) {
        console.error('Failed to fetch stats:', error);
      } finally {
        setLoading(false);
      }
    };

    fetchStats();
  }, []);

  return (
    <Layout>
      <div className="space-y-6">
        {/* Header */}
        <div>
          <h1 className="text-2xl font-bold">Dashboard</h1>
          <p className="text-muted-foreground">
            Hoş geldin, {user?.username}! Sunucu durumu ve istatistikler
          </p>
        </div>

        {loading ? (
          <div className="flex items-center justify-center h-64">
            <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600" />
          </div>
        ) : (
          <>
            {/* Stats Grid */}
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
              <StatCard
                title="Kullanıcılar"
                value={stats?.total_users || 0}
                icon={Users}
                color="bg-blue-500"
              />
              <StatCard
                title="Domainler"
                value={stats?.total_domains || 0}
                icon={Globe}
                color="bg-green-500"
              />
              <StatCard
                title="Veritabanları"
                value={stats?.total_databases || 0}
                icon={Database}
                color="bg-purple-500"
              />
              <StatCard
                title="E-posta Hesapları"
                value={stats?.total_emails || 0}
                icon={Mail}
                color="bg-orange-500"
              />
            </div>

            {/* System Stats */}
            {stats?.system_stats && (
              <>
                <h2 className="text-lg font-semibold">Sistem Kaynakları</h2>
                <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
                  <SystemStatCard
                    title="CPU Kullanımı"
                    used={stats.system_stats.cpu_usage}
                    total={100}
                    icon={Cpu}
                  />
                  <SystemStatCard
                    title="Bellek"
                    used={stats.system_stats.memory_used}
                    total={stats.system_stats.memory_total}
                    icon={MemoryStick}
                  />
                  <SystemStatCard
                    title="Disk"
                    used={stats.system_stats.disk_used}
                    total={stats.system_stats.disk_total}
                    icon={HardDrive}
                  />
                </div>
              </>
            )}

            {/* Quick Actions */}
            <h2 className="text-lg font-semibold">Hızlı İşlemler</h2>
            <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
              <Link to="/domains">
                <Button variant="outline" className="w-full h-auto py-4 flex-col gap-2">
                  <Globe className="w-5 h-5" />
                  <span>Domain Ekle</span>
                </Button>
              </Link>
              <Button variant="outline" className="h-auto py-4 flex-col gap-2" disabled>
                <Database className="w-5 h-5" />
                <span>Veritabanı Oluştur</span>
              </Button>
              <Button variant="outline" className="h-auto py-4 flex-col gap-2" disabled>
                <Mail className="w-5 h-5" />
                <span>E-posta Hesabı</span>
              </Button>
              <Button variant="outline" className="h-auto py-4 flex-col gap-2" disabled>
                <Users className="w-5 h-5" />
                <span>Kullanıcı Ekle</span>
              </Button>
            </div>
          </>
        )}
      </div>
    </Layout>
  );
}
