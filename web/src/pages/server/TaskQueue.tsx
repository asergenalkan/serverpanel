import { useState, useEffect } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/Card';
import { Button } from '@/components/ui/Button';
import { ListTodo, RefreshCw, Mail, Clock, AlertTriangle, CheckCircle } from 'lucide-react';
import Layout from '@/components/Layout';
import api from '@/lib/api';

interface MailQueueItem {
  id: string;
  sender: string;
  recipient: string;
  size: string;
  time: string;
  status: string;
}

interface CronJob {
  user: string;
  schedule: string;
  command: string;
  next_run: string;
}

interface QueueData {
  mail_queue: MailQueueItem[];
  mail_queue_count: number;
  cron_jobs: CronJob[];
  pending_tasks: number;
}

export default function TaskQueuePage() {
  const [data, setData] = useState<QueueData | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [activeTab, setActiveTab] = useState<'mail' | 'cron'>('mail');

  useEffect(() => {
    fetchQueueData();
    const interval = setInterval(fetchQueueData, 10000);
    return () => clearInterval(interval);
  }, []);

  const fetchQueueData = async () => {
    try {
      const response = await api.get('/server/queue');
      if (response.data.success) {
        setData(response.data.data);
      }
    } catch (err: any) {
      setError(err.response?.data?.error || 'Kuyruk bilgileri alınamadı');
    } finally {
      setLoading(false);
    }
  };

  const flushMailQueue = async () => {
    try {
      await api.post('/server/queue/flush');
      fetchQueueData();
    } catch (err: any) {
      setError(err.response?.data?.error || 'Kuyruk temizlenemedi');
    }
  };

  return (
    <Layout>
      <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold">Task Queue</h1>
          <p className="text-muted-foreground">Mail kuyruğu ve zamanlanmış görevler</p>
        </div>
        <Button onClick={fetchQueueData} variant="outline" size="sm">
          <RefreshCw className="w-4 h-4 mr-2" />
          Yenile
        </Button>
      </div>

      {/* Stats */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <Card>
          <CardContent className="pt-6">
            <div className="flex items-center gap-4">
              <div className="p-3 bg-blue-500/10 rounded-lg">
                <Mail className="w-6 h-6 text-blue-500" />
              </div>
              <div>
                <p className="text-2xl font-bold">{data?.mail_queue_count || 0}</p>
                <p className="text-sm text-muted-foreground">Mail Kuyruğunda</p>
              </div>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="pt-6">
            <div className="flex items-center gap-4">
              <div className="p-3 bg-green-500/10 rounded-lg">
                <Clock className="w-6 h-6 text-green-500" />
              </div>
              <div>
                <p className="text-2xl font-bold">{data?.cron_jobs?.length || 0}</p>
                <p className="text-sm text-muted-foreground">Cron Job</p>
              </div>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="pt-6">
            <div className="flex items-center gap-4">
              <div className="p-3 bg-orange-500/10 rounded-lg">
                <ListTodo className="w-6 h-6 text-orange-500" />
              </div>
              <div>
                <p className="text-2xl font-bold">{data?.pending_tasks || 0}</p>
                <p className="text-sm text-muted-foreground">Bekleyen Görev</p>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Tabs */}
      <div className="flex gap-2 border-b">
        <button
          onClick={() => setActiveTab('mail')}
          className={`px-4 py-2 font-medium transition-colors ${
            activeTab === 'mail'
              ? 'text-primary border-b-2 border-primary'
              : 'text-muted-foreground hover:text-foreground'
          }`}
        >
          <Mail className="w-4 h-4 inline mr-2" />
          Mail Kuyruğu
        </button>
        <button
          onClick={() => setActiveTab('cron')}
          className={`px-4 py-2 font-medium transition-colors ${
            activeTab === 'cron'
              ? 'text-primary border-b-2 border-primary'
              : 'text-muted-foreground hover:text-foreground'
          }`}
        >
          <Clock className="w-4 h-4 inline mr-2" />
          Cron Jobs
        </button>
      </div>

      {loading ? (
        <div className="flex items-center justify-center h-32">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
        </div>
      ) : error ? (
        <div className="bg-destructive/10 text-destructive p-4 rounded-lg">{error}</div>
      ) : (
        <>
          {/* Mail Queue Tab */}
          {activeTab === 'mail' && (
            <Card>
              <CardHeader className="flex flex-row items-center justify-between">
                <CardTitle className="flex items-center gap-2">
                  <Mail className="w-5 h-5" />
                  Mail Kuyruğu
                </CardTitle>
                {(data?.mail_queue_count || 0) > 0 && (
                  <Button variant="destructive" size="sm" onClick={flushMailQueue}>
                    Kuyruğu Temizle
                  </Button>
                )}
              </CardHeader>
              <CardContent>
                {!data?.mail_queue || data.mail_queue.length === 0 ? (
                  <div className="text-center py-8">
                    <CheckCircle className="w-12 h-12 text-green-500 mx-auto mb-4" />
                    <p className="text-muted-foreground">Mail kuyruğu boş</p>
                  </div>
                ) : (
                  <div className="overflow-x-auto">
                    <table className="w-full">
                      <thead>
                        <tr className="border-b">
                          <th className="text-left py-3 px-4 font-medium">ID</th>
                          <th className="text-left py-3 px-4 font-medium">Gönderen</th>
                          <th className="text-left py-3 px-4 font-medium">Alıcı</th>
                          <th className="text-right py-3 px-4 font-medium">Boyut</th>
                          <th className="text-left py-3 px-4 font-medium">Zaman</th>
                          <th className="text-left py-3 px-4 font-medium">Durum</th>
                        </tr>
                      </thead>
                      <tbody>
                        {data.mail_queue.map((item, index) => (
                          <tr key={index} className="border-b hover:bg-muted/50">
                            <td className="py-3 px-4 font-mono text-sm">{item.id}</td>
                            <td className="py-3 px-4 text-sm">{item.sender}</td>
                            <td className="py-3 px-4 text-sm">{item.recipient}</td>
                            <td className="py-3 px-4 text-right text-sm">{item.size}</td>
                            <td className="py-3 px-4 text-sm">{item.time}</td>
                            <td className="py-3 px-4">
                              {item.status === 'deferred' ? (
                                <span className="inline-flex items-center gap-1 text-yellow-500 text-sm">
                                  <AlertTriangle className="w-4 h-4" />
                                  Ertelendi
                                </span>
                              ) : (
                                <span className="text-sm">{item.status}</span>
                              )}
                            </td>
                          </tr>
                        ))}
                      </tbody>
                    </table>
                  </div>
                )}
              </CardContent>
            </Card>
          )}

          {/* Cron Jobs Tab */}
          {activeTab === 'cron' && (
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Clock className="w-5 h-5" />
                  Zamanlanmış Görevler (Cron)
                </CardTitle>
              </CardHeader>
              <CardContent>
                {!data?.cron_jobs || data.cron_jobs.length === 0 ? (
                  <div className="text-center py-8 text-muted-foreground">
                    Zamanlanmış görev bulunamadı.
                  </div>
                ) : (
                  <div className="overflow-x-auto">
                    <table className="w-full">
                      <thead>
                        <tr className="border-b">
                          <th className="text-left py-3 px-4 font-medium">Kullanıcı</th>
                          <th className="text-left py-3 px-4 font-medium">Zamanlama</th>
                          <th className="text-left py-3 px-4 font-medium">Komut</th>
                          <th className="text-left py-3 px-4 font-medium">Sonraki Çalışma</th>
                        </tr>
                      </thead>
                      <tbody>
                        {data.cron_jobs.map((job, index) => (
                          <tr key={index} className="border-b hover:bg-muted/50">
                            <td className="py-3 px-4 font-mono text-sm">{job.user}</td>
                            <td className="py-3 px-4 font-mono text-sm">{job.schedule}</td>
                            <td className="py-3 px-4 text-sm max-w-md truncate">
                              <code className="text-xs bg-muted px-2 py-1 rounded">{job.command}</code>
                            </td>
                            <td className="py-3 px-4 text-sm">{job.next_run}</td>
                          </tr>
                        ))}
                      </tbody>
                    </table>
                  </div>
                )}
              </CardContent>
            </Card>
          )}
        </>
      )}
    </div>
    </Layout>
  );
}
