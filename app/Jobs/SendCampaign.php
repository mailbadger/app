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
use Carbon\Carbon;

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
     * @var SesClient 
     */
    protected $client;

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
        $this->client = new SesClient($this->awsConfig); 
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
        $client = $this->client;

        $campaignService->updateCampaign(['status' => 'sending'], $this->campaign->id);

        $total = 0;

        $subscriberService->findSubscribersByListIdsByChunks($this->listIds, 1000, function ($subscribers) 
            use ($campaign, $emailService, $templateService, $client, &$total) {
            foreach($subscribers as $subscriber) {
                $token = $emailService->generateUniqueToken();

                $opensTrackerUrl = url('/api/emails/opens?_t='.$token);

                $tags = $this->createTagsFromSubscriberFields($subscriber->name, $subscriber->email, $subscriber->fields->toArray());

                $html = $templateService->renderTemplate($campaign->template_id, 
                    $subscriber->name, $subscriber->email, $opensTrackerUrl, $tags);

                $messageId = $emailService->sendEmail($client, $html, $subscriber->email,
                    $campaign->from_email, $campaign->from_name, $campaign->subject);

                $emailService->createSentEmail([
                    'subscriber_id' => $subscriber->id,
                    'campaign_id'   => $campaign->id,
                    'message_id'    => $messageId,
                    'token'         => $token,
                    'opens'         => 0,
                ]);

                $total++;
            }
        });

        $campaignService->updateCampaign(['status' => 'sent', 'recipients' => $total, 'sent_at' => Carbon::now()],
            $this->campaign->id);
    }

    /**
     * @param $subscriberName
     * @param $subscriberEmail
     * @param array $fields
     * @return array
     */
    private function createTagsFromSubscriberFields($subscriberName, $subscriberEmail, array $fields)
    { 
        $tags = [
            '/\*\|Name\|\*/i'  => $subscriberName,
            '/\*\|Email\|\*/i' => $subscriberEmail,
        ];

        foreach ($customFields as $key => $val) { 
            $tags['/\*\|' . $key . '\|\*/i'] = $val;
        } 

        return $tags;
    }
}
