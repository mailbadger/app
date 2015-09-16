<?php

namespace newsletters\Jobs;

use Illuminate\Contracts\Bus\SelfHandling;
use Illuminate\Contracts\Queue\ShouldQueue;
use Illuminate\Database\Eloquent\Collection;
use Illuminate\Queue\InteractsWithQueue;
use Illuminate\Queue\SerializesModels;
use newsletters\Entities\Campaign;
use newsletters\Services\EmailService;

class SendCampaign extends Job implements SelfHandling, ShouldQueue
{
    use InteractsWithQueue, SerializesModels;

    /**
     * @var Campaign
     */
    protected $campaign;

    /**
     * @var Collection
     */
    protected $subscribers;

    /**
     * Create a new job instance.
     *
     * @param Campaign $campaign
     * @param Collection $subscribers
     */
    public function __construct(Campaign $campaign, Collection $subscribers)
    {
        $this->campaign = $campaign;
        $this->subscribers = $subscribers;
    }

    /**
     * Execute the job.
     *
     * @param EmailService $emailService
     */
    public function handle(EmailService $emailService)
    {
        $campaign = $this->campaign;

        $this->subscribers->each(function ($subscriber) use ($campaign, $emailService) {
            $emailService->sendEmail($subscriber->email, $subscriber->name, $campaign->from_email, $campaign->from_name,
                $campaign->subject, $campaign->template_id, $subscriber->fields->toArray());

            $emailService->createSentEmail([
                'subscriber_id' => $subscriber->id,
                'campaign_id'   => $campaign->id,
                'opens'         => 0,
            ]);
        });
    }
}
