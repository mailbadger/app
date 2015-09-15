<?php

namespace newsletters\Jobs;

use Illuminate\Contracts\Bus\SelfHandling;
use Illuminate\Contracts\Queue\ShouldQueue;
use Illuminate\Database\Eloquent\Collection;
use Illuminate\Queue\InteractsWithQueue;
use Illuminate\Queue\SerializesModels;
use newsletters\Entities\Campaign;

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
     * @return void
     */
    public function handle()
    {
        $this->subscribers->each(function($subscriber) {

        });
    }
}
