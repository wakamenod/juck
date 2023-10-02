import java.util.concurrent.Executors;
import java.util.concurrent.ScheduledExecutorService;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.atomic.AtomicInteger;

public class ScheduledTaskExample {

    public static void main(String[] args) {
        ScheduledExecutorService executor = Executors.newScheduledThreadPool(1);

        long startTime = System.currentTimeMillis();
        AtomicInteger counter = new AtomicInteger(0);

        Runnable task = () -> {
            int taskRunCount = counter.incrementAndGet();

            long currentTime = System.currentTimeMillis();
            System.out.printf("Task %d started at %d ms\n", taskRunCount, currentTime - startTime);

            if (taskRunCount == 3 ||taskRunCount == 1) {
                try {
                    Thread.sleep(1200);  // Block for 1.2 seconds
                } catch (InterruptedException e) {
                    e.printStackTrace();
                }
            }

            currentTime = System.currentTimeMillis();
            System.out.printf("Task %d ended at %d ms\n", taskRunCount, currentTime - startTime);

            if (taskRunCount == 4) {
                executor.shutdown();
            }
        };

        executor.scheduleAtFixedRate(task, 2400, 500, TimeUnit.MILLISECONDS);
    }
}
