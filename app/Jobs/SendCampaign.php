<?php

namespace newsletters\Jobs;

use Illuminate\Contracts\Bus\SelfHandling;
use Illuminate\Contracts\Queue\ShouldQueue;
use Illuminate\Queue\InteractsWithQueue;
use Illuminate\Queue\SerializesModels;
use newsletters\Entities\Campaign;
use newsletters\Services\CampaignService;
use newsletters\Services\EmailService;
use newsletters\Services\TemplateService;
use newsletters\Services\SubscriberService;
use Aws\Ses\SesClient;

class SendCampaign extends Job implements SelfHandling, ShouldQueue
{
    use InteractsWithQueue, SerializesModels;

    /**
     * @var Campaign
     */
    protected $campaign;

    /**
     * @var array
     */
    protected $listIds;

    /**
     * @var array
     */
    protected $awsConfig;

    /**
     * Create a new job instance.
     *
     * @param Campaign $campaign
     * @param Collection $subscribers
     */
    public function __construct(Campaign $campaign, array $listIds, array $awsConfig)
    {
        $this->campaign = $campaign;
        $this->listIds = $listIds;
        $this->awsConfig = $awsConfig;
    }

    /**
     * Execute the job.
     *
     * @param EmailService $emailService
     * @param CampaignService $campaignService
     */
    public function handle(
        EmailService $emailService,
        CampaignService $campaignService,
        SubscriberService $subscriberService,
        TemplateService $templateService
    ) {
        $campaign = $this->campaign;

        $client = new SesClient($this->awsConfig);

        $campaignService->updateCampaign(['status' => 'sending'], $this->campaign->id);

        $subscriberService->findSubscribersByListIdsByChunks($this->listIds, 1000, function ($subscribers) use ($campaign, $emailService, $templateService, $client) {
            foreach($subscribers as $subscriber) {
                $html = $templateService->renderTemplate($campaign->template_id, 'Test Recipient', $subscriber->email, $subscriber->fields->toArray());
                $messageId = $emailService->sendEmail($client, $html, $subscriber->email, $campaign->from_email, $campaign->from_name, $campaign->subject);

                $emailService->createSentEmail([
                    'subscriber_id' => $subscriber->id,
                    'campaign_id'   => $campaign->id,
                    'message_id'    => $messageId,
                    'opens'         => 0,
                ]);
            }
        });

        $campaignService->updateCampaign(['status' => 'sent'], $this->campaign->id);
    }
}
