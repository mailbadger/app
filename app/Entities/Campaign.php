<?php

namespace newsletters\Entities;

use Illuminate\Database\Eloquent\Model;
use Prettus\Repository\Contracts\Transformable;
use Prettus\Repository\Traits\TransformableTrait;

class Campaign extends Model implements Transformable
{
    use TransformableTrait;

    protected $table = "campaigns";

    protected $fillable = [
        'name',
        'subject',
        'from_name',
        'from_email',
        'status',
    ];

    public function template()
    {
        return $this->belongsTo('newsletters\Entities\Template');
    }

    public function tags()
    {
        return $this->belongsToMany('newsletters\Entities\Tag', 'campaigns_tags', 'campaign_id',
            'tag_id')->withTimestamps();
    }
}
